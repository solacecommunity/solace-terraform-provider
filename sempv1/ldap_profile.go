package sempv1

import (
	"context"
)

type LdapProfileHelper struct {
	Context       context.Context
	Client        *SempV1Client
	Name          string
	Enabled       bool
	AdminDN       string
	AdminPassword string
	Host          string
	Index         int
	BaseDN        string
	SearchFilter  string
}

func (p *LdapProfileHelper) CreateLdapProfile() error {
	err := p.Client.CreateLdapProfile(p.Context, p.Name)
	if err != nil {
		if result, ok := err.(*ExecuteResult); ok {
			if result.ExecuteResult.ReasonCode == SEMP_ERR_ALREADY_EXISTS {
				return p.UpdateLdapProfile()
			}
			return err
		}
		return err
	}
	return p.UpdateLdapProfile()
}

func (p *LdapProfileHelper) ReadLdapProfile() error {
	result, err := p.Client.GetLdapProfile(p.Context, p.Name)
	if err != nil {
		return err
	}
	profile := result.Rpc.Show.LdapProfile.LdapProfile

	p.AdminDN = profile.AdminDn
	p.SearchFilter = profile.Search.Filter
	p.BaseDN = profile.Search.BaseDn
	p.Enabled = profile.Shutdown == "No"

	if len(profile.LdapServers.LdapServer) == 1 {
		server := profile.LdapServers.LdapServer[0]
		p.Host = server.LdapURI
		p.Index = server.Index
	}
	return nil
}

func (p *LdapProfileHelper) UpdateLdapProfile() error {
	err := p.Client.SetLdapProfileServer(p.Context, p.Name, p.Host, p.Index)
	if err != nil {
		return err
	}
	err = p.Client.SetLdapProfileSeachBaseDN(p.Context, p.Name, p.BaseDN)
	if err != nil {
		return err
	}
	err = p.Client.SetLdapProfileSeachFilter(p.Context, p.Name, p.SearchFilter)
	if err != nil {
		return err
	}
	err = p.Client.SetLdapProfileAdminDN(p.Context, p.Name, p.AdminDN, p.AdminPassword)
	if err != nil {
		return err
	}
	err = p.Client.EnableLdapProfile(p.Context, p.Name, p.Enabled)
	if err != nil {
		return err
	}
	return nil
}

func (p *LdapProfileHelper) DeleteLdapProfile() error {
	return p.Client.DeleteLdapProfile(p.Context, p.Name)
}
