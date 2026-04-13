package settings

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yumenaka/comigo/config"
)

// userInfoConfigXData 统一生成设置页的 Alpine 初始状态与交互逻辑。
func userInfoConfigXData() string {
	cfg := config.GetCfg()
	initialState := map[string]any{
		"loginProtection":   cfg.LoginProtection,
		"username":          cfg.Username,
		"current_password":  "",
		"password":          "",
		"ReEnterPassword":   "",
		"showPassword":      false,
		"enableOAuthLogin":  cfg.EnableOAuthLogin,
		"oauthProviderType": cfg.OAuthProviderTypeNormalized(),
		"oauthProviderName": cfg.OAuthProviderName,
		"oauthClientID":     cfg.OAuthClientID,
		"oauthClientSecret": cfg.OAuthClientSecret,
		"oauthAuthURL":      cfg.OAuthAuthURL,
		"oauthTokenURL":     cfg.OAuthTokenURL,
		"oauthUserInfoURL":  cfg.OAuthUserInfoURL,
		"oauthRedirectURL":  cfg.OAuthRedirectURL,
		"oauthScopesText":   strings.Join(cfg.OAuthScopes, " "),
		"isFormChanged":     false,
		"saving":            false,
	}
	payload, err := json.Marshal(initialState)
	if err != nil {
		return "{}"
	}

	return fmt.Sprintf(`Object.assign(%s, {
		providerTranslationKey() {
			switch (this.oauthProviderType) {
				case 'github':
					return 'OAuthProviderTypeGitHub';
				case 'google':
					return 'OAuthProviderTypeGoogle';
				case 'facebook':
					return 'OAuthProviderTypeFacebook';
				default:
					return 'OAuthProviderTypeOther';
			}
		},
		isCustomOAuthProvider() {
			return this.oauthProviderType === 'other';
		},
		handleLoginProtectionChange(event) {
			this.loginProtection = event.target.checked;
		},
		handleOAuthLoginChange(event) {
			this.enableOAuthLogin = event.target.checked;
		},
		init() {
			[
				'loginProtection',
				'username',
				'current_password',
				'password',
				'ReEnterPassword',
				'enableOAuthLogin',
				'oauthProviderType',
				'oauthProviderName',
				'oauthClientID',
				'oauthClientSecret',
				'oauthAuthURL',
				'oauthTokenURL',
				'oauthUserInfoURL',
				'oauthRedirectURL',
				'oauthScopesText',
			].forEach(name => {
				this.$watch(name, () => {
					this.isFormChanged = true;
				});
			});
		},
		checkFormData() {
			if (!this.isFormChanged) {
				return false;
			}
			if (this.loginProtection) {
				if (this.username.trim() === '') {
					showToast(i18next.t('PromptSetUsername'), 'error');
					return false;
				}
				if (this.password !== this.ReEnterPassword) {
					showToast(i18next.t('ErrPasswordMismatch'), 'error');
					return false;
				}
			}
			if (this.enableOAuthLogin) {
				const requiredValues = [this.oauthClientID, this.oauthClientSecret];
				if (requiredValues.some(value => value.trim() === '')) {
					showToast(i18next.t('PromptCompleteOAuthConfig'), 'error');
					return false;
				}
				if (this.isCustomOAuthProvider()) {
					const customProviderValues = [
						this.oauthAuthURL,
						this.oauthTokenURL,
						this.oauthUserInfoURL,
					];
					if (customProviderValues.some(value => value.trim() === '')) {
						showToast(i18next.t('PromptCompleteOAuthConfig'), 'error');
						return false;
					}
				}
			}
			return true;
		},
	})`, string(payload))
}
