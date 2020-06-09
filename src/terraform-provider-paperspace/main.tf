provider "paperspace" {
  region = "East Coast (NY2)"
  api_key = "12345678..."
}

data "paperspace_template" "my-template-1" {
  label = "Ubuntu 18.04 Server"
}

data "paperspace_user" "my-user-1" {
  email = "user@email.com"
  team_id = "te1234567"
}

resource "paperspace_script" "my-script-1" {
  name = "My Script"
  description = "a short description"
  script_text = <<EOF
#!/bin/bash
echo "Hello, World" > index.html
ufw allow 8080
nohup busybox httpd -f -p 8080 &
EOF
  is_enabled = true
  run_once = false
}

resource "paperspace_machine" "my-machine-1" {
  region = "East Coast (NY2)" // optional, defaults to provider region if not specified
  name = "Terraform Test"
  machine_type = "C1"
  size = 50
  billing_type = "hourly"
  assign_public_ip = true // optional, remove if you don't want a public ip assigned
  template_id = data.paperspace_template.my-template-1.id
  user_id = data.paperspace_user.my-user-1.id  // optional, remove to default
  team_id = data.paperspace_user.my-user-1.team_id
  script_id = paperspace_script.my-script-1.id // optional, remove for no script
}
