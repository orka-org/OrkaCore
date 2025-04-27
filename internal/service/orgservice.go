package service

import (
	"context"

	pb "github.com/orka-org/orkacore/api/organization/v1"
)

type OrgServiceService struct {
	pb.UnimplementedOrgServiceServer
}

func NewOrgServiceService() *OrgServiceService {
	return &OrgServiceService{}
}

func (s *OrgServiceService) CreateOrg(ctx context.Context, req *pb.CreateOrgRequest) (*pb.CreateOrgResponse, error) {
	return &pb.CreateOrgResponse{}, nil
}
func (s *OrgServiceService) GetOrg(ctx context.Context, req *pb.GetOrgRequest) (*pb.GetOrgResponse, error) {
	return &pb.GetOrgResponse{}, nil
}
func (s *OrgServiceService) UpdateOrg(ctx context.Context, req *pb.UpdateOrgRequest) (*pb.UpdateOrgResponse, error) {
	return &pb.UpdateOrgResponse{}, nil
}
func (s *OrgServiceService) DeleteOrg(ctx context.Context, req *pb.DeleteOrgRequest) (*pb.DeleteOrgResponse, error) {
	return &pb.DeleteOrgResponse{}, nil
}
func (s *OrgServiceService) UpdateOrgSettings(ctx context.Context, req *pb.UpdateOrgSettingsRequest) (*pb.UpdateOrgSettingsResponse, error) {
	return &pb.UpdateOrgSettingsResponse{}, nil
}
func (s *OrgServiceService) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (*pb.GetMembersResponse, error) {
	return &pb.GetMembersResponse{}, nil
}
func (s *OrgServiceService) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.InviteMemberResponse, error) {
	return &pb.InviteMemberResponse{}, nil
}
func (s *OrgServiceService) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	return &pb.RemoveMemberResponse{}, nil
}
func (s *OrgServiceService) UpdateMemberRole(ctx context.Context, req *pb.UpdateMemberRoleRequest) (*pb.UpdateMemberRoleResponse, error) {
	return &pb.UpdateMemberRoleResponse{}, nil
}
func (s *OrgServiceService) AddOrgRole(ctx context.Context, req *pb.AddOrgRoleRequest) (*pb.AddOrgRoleResponse, error) {
	return &pb.AddOrgRoleResponse{}, nil
}
func (s *OrgServiceService) RemoveOrgRole(ctx context.Context, req *pb.RemoveOrgRoleRequest) (*pb.RemoveOrgRoleResponse, error) {
	return &pb.RemoveOrgRoleResponse{}, nil
}
func (s *OrgServiceService) UpdateOrgRolePermission(ctx context.Context, req *pb.UpdateOrgRolePermissionRequest) (*pb.UpdateOrgRolePermissionResponse, error) {
	return &pb.UpdateOrgRolePermissionResponse{}, nil
}
