package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ldap.v2"
	"net/http"
)

const (
	ldapServer   = "192.168.3.102:389"
	ldapBindUser = "cn=admin,dc=ailieyun,dc=com" // LDAP 绑定用户
	ldapBindPass = "123456"                      // LDAP 绑定用户密码
	ldapBaseDN   = "dc=ailieyun,dc=com"          // LDAP 基本 DN
)

func LdapLogin(c *gin.Context) {
	var userloginform UserLoginForm
	if err := c.ShouldBindJSON(&userloginform); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数输入错误"})
		return
	}
	// 进行LDAP认证
	if err := authenticateLDAP(userloginform.UserName, userloginform.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

}
func authenticateLDAP(username, password string) error {
	l, err := ldap.Dial("tcp", ldapServer)
	if err != nil {
		return fmt.Errorf("ldap 链接失败 LDAP connection failed: %v", err)
	}
	defer l.Close()

	// 进行LDAP绑定
	err = l.Bind(ldapBindUser, ldapBindPass)
	if err != nil {
		return fmt.Errorf("LDAP bind failed: %v", err)
	}

	// 搜索用户
	searchRequest := ldap.NewSearchRequest(
		ldapBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(cn=%s))", username),
		[]string{"cn", "dn", "dc", "ou"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("LDAP search failed: %v", err)
	}

	if len(sr.Entries) != 1 {
		return fmt.Errorf("用户名没找到 User not found or multiple users found")
	}

	userDN := sr.Entries[0].DN

	// 尝试进行身份验证
	err = l.Bind(userDN, password)
	if err != nil {
		return fmt.Errorf("LDAP authentication failed: %s", "用户名或密码输入错误")
	}

	return nil
}
