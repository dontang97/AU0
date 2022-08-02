package dto

type Scope int

func (scope Scope) Int() int {
	return int(scope)
}

const (
	Issuetoken Scope = iota
	Email
	Birthday
	Id_Account_Query
	Pns_Message_Post
	Pns_Message_Query
	Profile_Write
	Pns_Message_Push
	Captcha_Verify
	Profile_Subscription_Write
	Payment_Security_Write
	Profile_Avatar_Write
	Profile_Advance_Write
	Enterprise_Profile_Read
	Enterprise_Profile_Write
	Addon_Arcade_Write
	Addon_Wl_Write
	Addon_Enterprisestore_Write
	Vrpns_Message_Post
	Vrpns_Message_Query
	Security_Write
)

func (scope Scope) String() string {
	switch scope {
	case Issuetoken:
		return "issuetoken"
	case Email:
		return "email"
	case Birthday:
		return "birthday"
	case Id_Account_Query:
		return "id.account.query"
	case Pns_Message_Post:
		return "pns.message.post"
	case Pns_Message_Query:
		return "pns.message.query"
	case Profile_Write:
		return "profile.write"
	case Pns_Message_Push:
		return "pns.message.push"
	case Captcha_Verify:
		return "captcha.verify"
	case Profile_Subscription_Write:
		return "profile.subscription.write"
	case Payment_Security_Write:
		return "payment.security.write"
	case Profile_Avatar_Write:
		return "profile.avatar.write"
	case Profile_Advance_Write:
		return "profile.advance.write"
	case Enterprise_Profile_Read:
		return "enterprise.profile.read"
	case Enterprise_Profile_Write:
		return "enterprise.profile.write"
	case Addon_Arcade_Write:
		return "addon.arcade.write"
	case Addon_Wl_Write:
		return "addon.wl.write"
	case Addon_Enterprisestore_Write:
		return "addon.enterprisestore.write"
	case Vrpns_Message_Post:
		return "vrpns.message.post"
	case Vrpns_Message_Query:
		return "vrpns.message.query"
	case Security_Write:
		return "security.write"
	default:
		return ""
	}
}

var scopeMap map[string]Scope = map[string]Scope{
	Issuetoken.String():                  Issuetoken,
	Email.String():                       Email,
	Birthday.String():                    Birthday,
	Id_Account_Query.String():            Id_Account_Query,
	Pns_Message_Post.String():            Pns_Message_Post,
	Pns_Message_Query.String():           Pns_Message_Query,
	Profile_Write.String():               Profile_Write,
	Pns_Message_Push.String():            Pns_Message_Push,
	Captcha_Verify.String():              Captcha_Verify,
	Profile_Subscription_Write.String():  Profile_Subscription_Write,
	Payment_Security_Write.String():      Payment_Security_Write,
	Profile_Avatar_Write.String():        Profile_Avatar_Write,
	Profile_Advance_Write.String():       Profile_Advance_Write,
	Enterprise_Profile_Read.String():     Enterprise_Profile_Read,
	Enterprise_Profile_Write.String():    Enterprise_Profile_Write,
	Addon_Arcade_Write.String():          Addon_Arcade_Write,
	Addon_Wl_Write.String():              Addon_Wl_Write,
	Addon_Enterprisestore_Write.String(): Addon_Enterprisestore_Write,
	Vrpns_Message_Post.String():          Vrpns_Message_Post,
	Vrpns_Message_Query.String():         Vrpns_Message_Query,
	Security_Write.String():              Security_Write,
}

func ScopeStrs2bitmap(scopeStrs []string) ([]byte, error) {
	size := ((len(scopeMap) - 1) >> 3) + 1
	bitmap := make([]byte, size)

	for _, str := range scopeStrs {
		scope, ok := scopeMap[str]
		if !ok {
			//return nil, fmt.Errorf("undefined scope: %s", str)
			continue
		}
		idx := size - (scope.Int() >> 3) - 1
		bitmap[idx] = (bitmap[idx] | (0x1 << (scopeMap[str].Int() & 0x7)))
	}

	return bitmap, nil
}

func Bitmap2ScopeStrs(bitmap []byte) []string {
	scopes := []string{}
	size := len(bitmap)

	idx := 0
	for i := size - 1; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			scope := Scope(idx).String()
			if scope == "" { // index exceeds defined scopes
				return scopes
			}

			if (bitmap[i]>>j)&0x1 == 1 {
				scopes = append(scopes, scope)
			}
			idx += 1
		}
	}

	return scopes
}
