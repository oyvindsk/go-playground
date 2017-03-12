
https://developers.google.com/sheets/quickstart/go

Authorization when running locally, using Application Deffault Credentials:
(this shoud wok, but it does not ATM)
    gcloud auth application-default login --scopes="https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/spreadsheets.readonly"