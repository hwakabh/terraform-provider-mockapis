terraform {
  required_providers {
    mockapis = {
      source = "registry.terraform.io/hwakabh/mockapis"
    }
  }
}

provider "mockapis" {
  username = "hwakabh"
}
