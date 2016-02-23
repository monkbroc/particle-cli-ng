Releasing Toolbelt v4
=====================

Prerequisites:

* `gem install rake aws-sdk`
* `particle-cli-ng` repo
* `PARTICLE_CLI_RELEASE_ACCESS` and `PARTICLE_CLI_RELEASE_SECRET` for the S3 repository

Instructions:

* Increment version
* Update CHANGELOG
* [optional] Update dev branch with latest changes.
* [optional] Run `rake release` on dev branch to test changes.
* Run `rake release` on master branch.

Notes:

Make sure the user with the S3 credentials above have these permissions.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::particle-cli-ng-alpha"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:DeleteObject",
                "s3:GetObject",
                "s3:PutObject",
                "s3:PutObjectAcl"
            ],
            "Resource": [
                "arn:aws:s3:::particle-cli-ng-alpha/*"
            ]
        }
    ]
}
```
