// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"math"
	"time"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendor models.VendorVerifyInfo //need use a appvendor object to call function
	vendorPbSet := vendor.ParseVendorSet2PbSet(ctx, vendors)

	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: vendorPbSet,
		TotalCount:          vendorCount,
	}
	return res, nil
}

func (s *Server) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	err := s.validateSubmitParams(ctx, req)
	if err != nil {
		return nil, err
	}

	var userID string
	ifExist, err := s.checkVendorVerifyInfoIfExit(ctx, req.UserId)
	if err != nil {
		logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], %+v", req.UserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if ifExist {
		attributes := manager.BuildUpdateAttributes(req, constants.ColumnCompanyName, constants.ColumnCompanyWebsite, constants.ColumnCompanyProfile,
			constants.ColumnAuthorizerName, constants.ColumnAuthorizerEmail, constants.ColumnAuthorizerPhone, constants.ColumnBankName, constants.ColumnBankAccountName,
			constants.ColumnBankAccountNumber)
		attributes[constants.ColumnStatus] = constants.StatusSubmitted
		attributes[constants.ColumnSubmitTime] = time.Now()
		attributes[constants.ColumnStatusTime] = time.Now()
		logger.Debug(ctx, "SubmitVendorVerifyInfo got attributes: [%+v]", attributes)
		userID, err = UpdateVendorVerifyInfo(ctx, req.UserId, attributes)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%s], %+v", userID, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
		}
	} else {
		vendor := &models.VendorVerifyInfo{}
		vendor = vendor.ParseReq2Vendor(ctx, req)
		userID, err = CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%+v], %+v", vendor, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
		logger.Debug(ctx, "vendorVerifyInfo does not exit, create new vendorVerifyInfo verify info, [%+v]", vendor)

	}
	res := &pb.SubmitVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, nil
}

func (s *Server) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	appVendorUserID := req.GetUserId()
	sender := ctxutil.GetSender(ctx)
	approver := string(sender.UserId)
	userID, err := PassVendorVerifyInfo(ctx, appVendorUserID, approver)
	if err != nil {
		logger.Error(ctx, "Failed to pass vendorVerifyInfo [%s], %+v", userID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.PassVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (s *Server) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	appVendorUserID := req.GetUserId()
	sender := ctxutil.GetSender(ctx)
	approver := string(sender.UserId)
	rejectMsg := req.GetRejectMessage().GetValue()
	userID, err := RejectVendorVerifyInfo(ctx, appVendorUserID, rejectMsg, approver)
	if err != nil {
		logger.Error(ctx, "Failed to reject vendorVerifyInfo [%s], %+v", userID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (s *Server) DescribeAppVendorStatistics(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorStatisticsResponse, error) {

	appClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendorStatisticses []*models.VendorStatistics
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	/*============================================================================================================*/
	//To get ClusterCountTotal
	var clusterCntAll4AllPages int32
	var clusterCntMonth4AllPages int32
	for _, vendor := range vendors {
		//step1:Get real appCnt for each vendor
		var vendorStatistics models.VendorStatistics
		pbApps, appCnt, err := appClient.DescribeActiveAppsByMaxRow(ctx, db.DefaultSelectLimit, 0)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		//step2:if the real appCnt is smaller than db.DefaultSelectLimit,there is only one page apps,and the rows of this one page is length of pbApps.
		//Just accumulate the clusterCnt4SingleApp for each app.
		if appCnt <= int32(db.DefaultSelectLimit) {
			for _, pbApp := range pbApps {
				_, clusterCntAll4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), false, db.DefaultSelectLimit, 0)
				_, clusterCntMonth4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), true, db.DefaultSelectLimit, 0)
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
				}
				clusterCntAll4AllPages = clusterCntAll4AllPages + clusterCntAll4SingleApp
				clusterCntMonth4AllPages = clusterCntMonth4AllPages + clusterCntMonth4SingleApp
			}

		} else {
			//step3:if the real appCnt is bigger than db.DefaultSelectLimit(200),there are more than 1 page Apps.
			//Should accumulate the clusterCnt4SingleApp for each apps ,then accumulate the number for each page.
			pages := int(math.Ceil(float64(appCnt / db.DefaultSelectLimit)))
			for i := 1; i <= pages; i++ {
				offset := db.DefaultSelectLimit * i
				pbApps, _, err := appClient.DescribeActiveAppsByMaxRow(ctx, db.DefaultSelectLimit, uint32(offset))
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
				}

				var clusterCntAll4OnePage int32
				var clusterCntMonth4OnePage int32
				for _, pbApp := range pbApps {
					_, clusterCntAll4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), false, db.DefaultSelectLimit, uint32(offset))
					_, clusterCntMonth4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), true, db.DefaultSelectLimit, uint32(offset))
					if err != nil {
						return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
					}
					clusterCntAll4OnePage = clusterCntAll4OnePage + clusterCntAll4SingleApp
					clusterCntMonth4OnePage = clusterCntMonth4OnePage + clusterCntMonth4SingleApp
				}
				clusterCntAll4AllPages = clusterCntAll4AllPages + clusterCntAll4OnePage
				clusterCntMonth4AllPages = clusterCntMonth4AllPages + clusterCntMonth4OnePage
			}

		}

		/*============================================================================================================*/
		vendorStatistics.UserId = vendor.UserId

		vendorStatistics.CompanyName = vendor.CompanyName
		vendorStatistics.ActiveAppCount = int32(appCnt)
		vendorStatistics.ClusterCountTotal = clusterCntAll4AllPages
		vendorStatistics.ClusterCountMonth = clusterCntMonth4AllPages

		vendorStatisticses = append(vendorStatisticses, &vendorStatistics)

	}

	var vendorStatistics models.VendorStatistics //need use a vendorStatistics object to call function
	vendorStatisticsPbSet := vendorStatistics.ParseVendorStatisticsSet2PbSet(ctx, vendorStatisticses)

	res := &pb.DescribeVendorStatisticsResponse{
		VendorVerifyStatisticsSet: vendorStatisticsPbSet,
		TotalCount:                vendorCount,
	}
	return res, nil
}

func (s *Server) checkVendorVerifyInfoIfExit(ctx context.Context, userID string) (bool, error) {
	info, err := GetVendorVerifyInfo(ctx, userID)
	if info == nil && err != nil {
		return false, nil
	} else if info != nil && err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (s *Server) validateSubmitParams(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) error {
	url := req.CompanyWebsite.GetValue()
	isUrlFmt, err := VerifyUrl(ctx, url)

	if !isUrlFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.CompanyWebsite.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalUrlFormat, url)
	}

	email := req.AuthorizerEmail.GetValue()
	isEmailFmt, err := VerifyEmailFmt(ctx, email)

	if !isEmailFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerEmail.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalEmailFormat, email)
	}

	phone := req.AuthorizerPhone.GetValue()
	isPhoneFmt, err := VerifyPhoneFmt(ctx, phone)
	if !isPhoneFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerPhone.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalPhoneFormat, phone)
	}

	isBankAccountNumberFmt, err := VerifyBankAccountNumberFmt(ctx, req.BankAccountNumber.GetValue())
	if !isBankAccountNumberFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.BankAccountNumber.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalBankAccountNumberFormat, req.BankAccountNumber.GetValue())
	}
	return nil
}
