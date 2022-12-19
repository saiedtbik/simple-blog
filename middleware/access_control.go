package middleware

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/harranali/authority"
	"withCasbin/component"
	"withCasbin/handler"

	//_ "github.com/go-sql-driver/mysql"
	"log"
	"withCasbin/models"
)

// Authenticate determines if current subject has logged in.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session id
		sessionId, _ := c.Cookie("current_subject")
		// Get current subject
		sub, err := component.GlobalCache.Get(sessionId)
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			c.AbortWithStatusJSON(401, component.RestResponse{Message: "user hasn't logged in yet"})
			return
		}
		c.Set("current_subject", string(sub))
		c.Next()
	}
}

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, existed := c.Get("current_subject")
		if !existed {
			c.AbortWithStatusJSON(401, component.RestResponse{Message: "user hasn't logged in yet"})
			return
		}
		// casbin enforces policy
		var roleName =handler.FindUserRole(val.(string))
		ok, err := enforce(roleName, obj, act, adapter)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, component.RestResponse{Message: "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, component.RestResponse{Message: "forbidden"})
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	enforcer = LoadPolicyFromDB(enforcer)
//	enforcer.AddRoleForUser("bob","admin")
	//enforcer.AddPolicy("p", "admin", "post", "read")

	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}

func LoadPolicyFromDB(enforcer *casbin.Enforcer) (*casbin.Enforcer){
	var userRoles []models.UserRole
	models.DB.Find(&userRoles)

	var rolePermissionResources []models.RolePermissionResource
	models.DB.Find(&rolePermissionResources)

	for _,userRole := range userRoles{
		enforcer.AddRoleForUser(userRole.User,userRole.Role)
	}

	for _,rolePermRsrs := range rolePermissionResources{
		enforcer.AddPolicy(rolePermRsrs.Role,rolePermRsrs.Resource, rolePermRsrs.Permission)
	}
	return enforcer
}

func enforce2(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	// initiate authority
	auth := authority.New(authority.Options{
		TablesPrefix: "authority_",
		DB:           models.DB,
	})
	err := auth.CreateRole("admin")

	//	enforcer.AddRoleForUser("bob","admin")
	//enforcer.AddPolicy("p", "admin", "post", "read")

	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	err = auth.CreatePermission("read")
	err = auth.CreatePermission("write")

	err = auth.AssignPermissions("admin", []string{
		"read",
		"write",
	})
	err = auth.AssignRole(1, "admin") // 1 is the user id
	ok, err := auth.CheckRole(1, sub) //1 is the user id
	ok, err = auth.CheckPermission(1, act) // 1 is the user id
	ok, err = auth.CheckRolePermission(sub, act)
	//err = auth.RevokeRole(1, "role-a")
	return ok, err
}





