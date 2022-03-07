package manager_test

import (
	manager "github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jwt", func() {
	var (
		privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgER92Clgmc3ikcLbSjZFo2jC4mA+aP/kNKRKMud3AlDY0mPVbEu3
LeHov93zmvy1s5k5XTBPeAdRybKODQ0/jHOeXOOflaynDKZWknD8/WmU0O64Z/Qf
IH7c1FhDYX1VUZhwPwpL0IxYJDIoCKzwBafPsIC4PUH+Lqyga3emP/v1AgMBAAEC
gYA/wvQAmTi2Da3qvCFbYvscZQk/1foD9xv2skivaQBT6XX7kM1Ps4lYXUh5RPbN
Og6nn1qcxe6UydQ+kLWf1sBWT8xJP34RNm93dkHzSteU8WdlVYGNQqQQxYQXaSpN
g8kXjMY8+EUaQkptdQTgcpT2ZCW0ZD0LSpmklsRPPSW0xQJBAIabP/RAhZRMY6Eu
I1JDkSeJSnp7QP/M0HOu6tBcKxVjApaX3RUxIxR8e7F4TEUgaMmsU5TIfkKtdnDS
wLmcHG8CQQCCQpOyN4WffT0marIbSIJViaQYbyQAW/qrgYZrHNMApmv3dklLefF1
nifIHHQn6IkzBZaN0EfRlp907lu9bWfbAkB8SzdO75VpTvBgkR4EhGewvlGLr+xh
SFrjt40UQUd3RCnLrQd03h6qeBgv1Al5e2fHcdzr8gbEwzAvFizoN4L5AkAwzA4W
WlRVbg5FYPz92Yjx0FFH0gLTm6FpNGmNoMuu16lkl8xXWQRKgof2oCondSZIldRT
pe3xpxJvNIfri5u3AkA5VTxp77yXQ7Bra/F3eyLzQ8VzhhjXes+jxb6imag2Ry9o
AuBDWd8zTFaIkV0Wl8BteGrMMfhLv0F9JxuDcZas
-----END RSA PRIVATE KEY-----`
		publicKey = `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgER92Clgmc3ikcLbSjZFo2jC4mA+
aP/kNKRKMud3AlDY0mPVbEu3LeHov93zmvy1s5k5XTBPeAdRybKODQ0/jHOeXOOf
laynDKZWknD8/WmU0O64Z/QfIH7c1FhDYX1VUZhwPwpL0IxYJDIoCKzwBafPsIC4
PUH+Lqyga3emP/v1AgMBAAE=
-----END PUBLIC KEY-----`
		token     string
		provider  *manager.PublicKeySimpleProvider
		parser    *manager.RsaJwtTokenParser
		jwtClaims manager.ApigatewayJwtClaims
	)

	BeforeEach(func() {
		jwtClaims = manager.ApigatewayJwtClaims{
			App: &manager.ApigatewayJwtApp{
				AppCode:  "app_code",
				Verified: true,
			},
			User: &manager.ApigatewayJwtUser{
				Username:   "username",
				SourceType: "default",
				Verified:   true,
			},
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
		Expect(err).To(BeNil())

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
		jwtToken.Header["kid"] = "testing"
		token, err = jwtToken.SignedString(key)
		Expect(err).To(BeNil())

		provider = manager.NewPublicKeySimpleProvider(map[string]string{
			"testing": publicKey,
		})
		parser = manager.NewRsaJwtTokenParser(provider)
	})

	It("should parse token", func() {
		claims, err := parser.Parse(token)
		Expect(err).To(BeNil())

		Expect(claims.ApiName).To(Equal("testing"))
		Expect(claims.App).To(Equal(jwtClaims.App))
		Expect(claims.User).To(Equal(jwtClaims.User))
	})
})
