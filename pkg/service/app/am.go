// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateAppRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("name", "version_package", "version_type").
			StringChosen("version_type", repoiface.SupportedPackageType).
			Exec()
	case *pb.ValidatePackageRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_package", "version_type").
			StringChosen("version_type", repoiface.SupportedPackageType).
			Exec()
	case *pb.ModifyAppRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("app_id").
			Exec()
	case *pb.DeleteAppsRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("app_id").
			Exec()
	case *pb.CreateAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("app_id", "package", "type").
			StringChosen("type", repoiface.SupportedPackageType).
			Exec()
	case *pb.ModifyAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id").
			Exec()
	case *pb.GetAppVersionPackageRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	case *pb.GetAppVersionPackageFilesRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	case *pb.SubmitAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id").
			Exec()
	case *pb.ReviewAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id", "role").
			StringChosen("role", reviewSupportRoles).
			Exec()
	case *pb.CancelAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id").
			Exec()
	case *pb.ReleaseAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id").
			Exec()
	case *pb.DeleteAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllDeveloperRoles).
			Required("version_id").
			Exec()
	case *pb.PassAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("version_id", "role").
			StringChosen("role", reviewSupportRoles).
			Exec()
	case *pb.RejectAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("version_id", "message", "role").
			StringChosen("role", reviewSupportRoles).
			Exec()
	case *pb.SuspendAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("version_id").
			Exec()
	case *pb.RecoverAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("version_id").
			Exec()
	case *pb.GetAppStatisticsRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Exec()
	}
	return nil
}

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	sender := ctxutil.GetSender(ctx)
	switch r := req.(type) {
	case *pb.DescribeAppsRequest:
		if !sender.IsGlobalAdmin() && !sender.IsDeveloper() {
			r.Status = []string{constants.StatusActive}
		}
		return r
	case *pb.DescribeAppVersionsRequest:
		if !sender.IsGlobalAdmin() && !sender.IsDeveloper() {
			r.Status = []string{constants.StatusActive}
		}
		return r
	}
	return req
}
