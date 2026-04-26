package settings

import (
	"encoding/json"
	"fmt"

	"github.com/yumenaka/comigo/config"
)

// userInfoConfigXData 统一生成设置页的 Alpine 初始状态与交互逻辑。
func userInfoConfigXData() string {
	cfg := config.GetCfg()
	initialState := map[string]any{
		"username":         cfg.Username,
		"current_password": "",
		"password":         "",
		"ReEnterPassword":  "",
		"showPassword":     false,
		"isFormChanged":    false,
		"saving":           false,
	}
	payload, err := json.Marshal(initialState)
	if err != nil {
		return "{}"
	}

	return fmt.Sprintf(`Object.assign(%s, {
		init() {
			[
				'username',
				'current_password',
				'password',
				'ReEnterPassword',
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
			if (this.password !== this.ReEnterPassword) {
				showToast(i18next.t('err_password_mismatch'), 'error');
				return false;
			}
			if (this.username.trim() === '' && (this.password !== '' || this.ReEnterPassword !== '')) {
				showToast(i18next.t('prompt_set_username'), 'error');
				return false;
			}
			return true;
		},
	})`, string(payload))
}
