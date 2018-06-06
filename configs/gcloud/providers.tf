provider "google" {
  credentials = "${file("google-cred.json")}"
  project     = "${var.gcloud-project}"
  region      = "${var.gcloud-region}"
}
