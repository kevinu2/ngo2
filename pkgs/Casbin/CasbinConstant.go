package Casbin

const (
	DefaultClaims = "claims"
	DefaultPrefix = "/"
	Admin         = "ALL"
	AdminSub      = "1"
	ObjOffset     = "?"
	Content       = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && (r.act == p.act|| p.act == "*") ||  r.sub == "1"`
)
