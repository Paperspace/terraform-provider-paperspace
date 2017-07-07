provider "paperspace" {
  apiKey = "1be4f97..."
  region = "East Coast (NY2)"
}

resource "paperspace_machine" "my-machine-1" {
  region = "East Coast (NY2)" // or "${var.paperspace_region}"
  machineType = "C1"
  size = 50
  billingType = "hourly"
  machineName = "Tom Terraform Test 1",
  templateId = "tqalmii" // Ubuntu 16.04 Server
  //userId = "hoo"
}
