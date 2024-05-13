package webhooks

var (
	AvatarURL = "https://cdn.discordapp.com/attachments/793996798305239071/1232550389790801991/529_Branding_Full_Export.png?ex=6629dd80&is=66288c00&hm=9ef0b5c5f91c94c1e2ba4503bb44bf366e4580365b32461c0434d362650e337d&"
)

func StrPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}
