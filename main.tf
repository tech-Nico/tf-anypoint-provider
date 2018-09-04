provider "anypoint" {
    username = "the-username"
    password = "the-password"
}

resource "anypoint_business_group" "Org1" {
  path_to_bg = "org2"
}