resource "jamfpro_local_admin_password_settings" "local_admin_password_settings_001" {
  auto_deploy_enabled                 = false
  password_rotation_time_seconds      = 3600
  auto_rotate_enabled                 = false
  auto_rotate_expiration_time_seconds = 7776000
}
