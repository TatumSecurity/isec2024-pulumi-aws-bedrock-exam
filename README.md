<h1 align="center">
isec2024-pulumi-aws-bedrock-exam
</h1>

<p align="center">
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"/></a>
<a href="https://github.com/Tatumsecurity/isec2024-pulumi-aws-bedrock-exam/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"/></a>
<a href="https://github.com/Tatumsecurity/isec2024-pulumi-aws-bedrock-exam/releases/latest"><img src="https://img.shields.io/github/v/release/Tatumsecurity/isec2024-pulumi-aws-bedrock-exam" /></a>
</p>
<p align="center">
<br>

### 참고) 해당 리포는 S3가 취약으로 설정되어있습니다. 실제 테스트 목적으로 사용할 경우 Pulumi.pulumi-cloud-ai-exam.yaml 내 PublicAccess를 false로 설정 부탁드립니다.

## 설치

### 필수 패키지

- [Pulumi](https://www.pulumi.com/docs/install)
- [Golang](https://go.dev/doc/install)
- [AWS CLI](https://docs.aws.amazon.com/ko_kr/cli/latest/userguide/getting-started-install.html)

## 종속성

```bash
go mod download
```

### 인증

해당 프로젝트는 pulumi를 통해서 인프라 배포를 진행합니다. 그에 따라서 [해당 문서](https://www.pulumi.com/registry/packages/aws/installation-configuration/)와 같은 AWS인증이 필요합니다.

#### AWS CLI

host 내 ak / sk를 입력하는 방법으로 진행할 경우 사용합니다.

```bash
export AWS_ACCESS_KEY_ID="ak 입력"
export AWS_SECRET_ACCESS_KEY="sk 입력"
```

#### AWS profile

host 내 aws profile을 입력하는 방법으로 진행할 경우 사용합니다.

```ini
[isec2024]
aws_access_key_id=ak 입력
aws_secret_access_key=sk 입력
```

### backend

pulumi도 terraform과 유사하게 인프라에 대한 상태관리를 수행합니다. 자세한 백엔드 관리는 [해당문서](https://www.pulumi.com/docs/iac/concepts/state-and-backends/) 참고 부탁드립니다. 해당 프로젝트는 local을 통해서 메타데이터를 관리합니다.

#### 에제

##### local

```bash
$ pulumi login --local
```

## stack

pulumi는 스택이라는 부분을 사용해서 프로젝트 구분을 진행합니다. 스택은 환경별로 독립적으로 인프라 배포를 진행할 수 있으며 일반적으로 dev, stage, prod 등 환경으로 분리를 하거나, 기능별로 분리합니다.

### stack 생성

`pulumi stack init 스택이름`을 통해서 스택을 초기화 할 수 있습니다.

#### 예제

```bash
$ pulumi stack init pulumi-cloud-ai-exam.yaml
```

## 인프라 배포

```bash
$ pulumi up
Previewing update (pulumi-cloud-ai-exam):
     Type                               Name                                                    Plan
 +   pulumi:pulumi:Stack                pulumi-cloud-ai-exam-pulumi-cloud-ai-exam               create
 +   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-input-s3                          create
 +   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-output-s3                         create
 +   ├─ aws:iam:Role                    TatumImageSummarizationLambdaRole                       create
 +   ├─ aws:s3:BucketPublicAccessBlock  tatumsecurity-bedrock-output-s3PublicAccessBlock        create
 +   ├─ aws:cloudwatch:EventRule        TatumInputBucketRule                                    create
 +   ├─ aws:lambda:Function             TatumImageSummarizationFilterLabels                     create
 +   ├─ aws:lambda:Function             TatumImageSummarizationBuildOutput                      create
 +   ├─ aws:s3:BucketNotification       tatumsecurity-bedrock-input-s3Notification              create
 +   ├─ aws:s3:BucketOwnershipControls  tatumsecurity-bedrock-output-s3BucketOwnershipControls  create
 +   ├─ aws:s3:BucketPolicy             tatumsecurity-bedrock-output-s3Policy                   create
 +   ├─ aws:iam:Role                    TatumImageSummarizationSfnRole                          create
 +   ├─ aws:sfn:StateMachine            TatumImageSummarizationStateMachine                     create
 +   ├─ aws:s3:BucketAclV2              tatumsecurity-bedrock-output-s3Acls                     create
 +   ├─ aws:iam:Policy                  TatumImageSummarizationRulePolicy                       create
 +   ├─ aws:iam:Policy                  TatumImageSummarizationSfnPolicy                        create
 +   ├─ aws:iam:Role                    TatumImageSummarizationRuleRole                         create
 +   └─ aws:cloudwatch:EventTarget      TatumInputRuleTarget                                    create

Resources:
    + 18 to create

Do you want to perform this update?  [Use arrows to move, type to filter]
  yes
> no
  details
```

```
Resources:
Do you want to perform this update? yes
Updating (pulumi-cloud-ai-exam):
     Type                               Name                                                    Status
 +   pulumi:pulumi:Stack                pulumi-cloud-ai-exam-pulumi-cloud-ai-exam               created (29s)
 +   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-input-s3                          created (3s)
 +   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-output-s3                         created (3s)
 +   ├─ aws:iam:Role                    TatumImageSummarizationLambdaRole                       created (1s)
 +   ├─ aws:cloudwatch:EventRule        TatumInputBucketRule                                    created (1s)
 +   ├─ aws:lambda:Function             TatumImageSummarizationFilterLabels                     created (16s)
 +   ├─ aws:lambda:Function             TatumImageSummarizationBuildOutput                      created (16s)
 +   ├─ aws:s3:BucketPublicAccessBlock  tatumsecurity-bedrock-output-s3PublicAccessBlock        created (1s)
 +   ├─ aws:s3:BucketOwnershipControls  tatumsecurity-bedrock-output-s3BucketOwnershipControls  created (1s)
 +   ├─ aws:s3:BucketNotification       tatumsecurity-bedrock-input-s3Notification              created (1s)
 +   ├─ aws:s3:BucketAclV2              tatumsecurity-bedrock-output-s3Acls                     created (0.65s)
 +   ├─ aws:s3:BucketPolicy             tatumsecurity-bedrock-output-s3Policy                   created (1s)
 +   ├─ aws:iam:Policy                  TatumImageSummarizationSfnPolicy                        created (0.78s)
 +   ├─ aws:iam:Role                    TatumImageSummarizationSfnRole                          created (1s)
 +   ├─ aws:sfn:StateMachine            TatumImageSummarizationStateMachine                     created (2s)
 +   ├─ aws:iam:Policy                  TatumImageSummarizationRulePolicy                       created (0.78s)
 +   ├─ aws:iam:Role                    TatumImageSummarizationRuleRole                         created (1s)
 +   └─ aws:cloudwatch:EventTarget      TatumInputRuleTarget                                    created (1s)

Resources:
    + 18 created

Duration: 34s
```

## 인프라 삭제

스택까지 완벽하게 삭제하고싶을 경우 실행합니다.

```bash
$ pulumi destroy
Previewing destroy (pulumi-cloud-ai-exam):
     Type                               Name                                                    Plan
 -   pulumi:pulumi:Stack                pulumi-cloud-ai-exam-pulumi-cloud-ai-exam               delete
 -   ├─ aws:s3:BucketNotification       tatumsecurity-bedrock-input-s3Notification              delete
 -   ├─ aws:lambda:Function             TatumImageSummarizationBuildOutput                      delete
 -   ├─ aws:iam:Role                    TatumImageSummarizationLambdaRole                       delete
 -   ├─ aws:s3:BucketPolicy             tatumsecurity-bedrock-output-s3Policy                   delete
 -   ├─ aws:s3:BucketAclV2              tatumsecurity-bedrock-output-s3Acls                     delete
 -   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-input-s3                          delete
 -   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-output-s3                         delete
 -   ├─ aws:iam:Role                    TatumImageSummarizationSfnRole                          delete
 -   ├─ aws:s3:BucketOwnershipControls  tatumsecurity-bedrock-output-s3BucketOwnershipControls  delete
 -   ├─ aws:iam:Policy                  TatumImageSummarizationSfnPolicy                        delete
 -   ├─ aws:cloudwatch:EventRule        TatumInputBucketRule                                    delete
 -   ├─ aws:iam:Role                    TatumImageSummarizationRuleRole                         delete
 -   ├─ aws:sfn:StateMachine            TatumImageSummarizationStateMachine                     delete
 -   ├─ aws:s3:BucketPublicAccessBlock  tatumsecurity-bedrock-output-s3PublicAccessBlock        delete
 -   ├─ aws:lambda:Function             TatumImageSummarizationFilterLabels                     delete
 -   ├─ aws:cloudwatch:EventTarget      TatumInputRuleTarget                                    delete
 -   └─ aws:iam:Policy                  TatumImageSummarizationRulePolicy                       delete

Resources:
    - 18 to delete

Do you want to perform this destroy?  [Use arrows to move, type to filter]
  yes
> no
  details
```

```
Destroying (pulumi-cloud-ai-exam):
     Type                               Name                                                    Status
 -   pulumi:pulumi:Stack                pulumi-cloud-ai-exam-pulumi-cloud-ai-exam               deleted (0.00s)
 -   ├─ aws:cloudwatch:EventTarget      TatumInputRuleTarget                                    deleted (2s)
 -   ├─ aws:iam:Role                    TatumImageSummarizationRuleRole                         deleted (1s)
 -   ├─ aws:iam:Policy                  TatumImageSummarizationRulePolicy                       deleted (0.61s)
 -   ├─ aws:sfn:StateMachine            TatumImageSummarizationStateMachine                     deleted (56s)
 -   ├─ aws:s3:BucketPolicy             tatumsecurity-bedrock-output-s3Policy                   deleted (0.91s)
 -   ├─ aws:iam:Role                    TatumImageSummarizationSfnRole                          deleted (0.97s)
 -   ├─ aws:iam:Policy                  TatumImageSummarizationSfnPolicy                        deleted (0.55s)
 -   ├─ aws:s3:BucketAclV2              tatumsecurity-bedrock-output-s3Acls                     deleted (0.01s)
 -   ├─ aws:s3:BucketPublicAccessBlock  tatumsecurity-bedrock-output-s3PublicAccessBlock        deleted (1s)
 -   ├─ aws:s3:BucketNotification       tatumsecurity-bedrock-input-s3Notification              deleted (0.74s)
 -   ├─ aws:s3:BucketOwnershipControls  tatumsecurity-bedrock-output-s3BucketOwnershipControls  deleted (0.97s)
 -   ├─ aws:lambda:Function             TatumImageSummarizationFilterLabels                     deleted (0.94s)
 -   ├─ aws:lambda:Function             TatumImageSummarizationBuildOutput                      deleted (0.92s)
 -   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-output-s3                         deleted (0.73s)
 -   ├─ aws:iam:Role                    TatumImageSummarizationLambdaRole                       deleted (0.49s)
 -   ├─ aws:s3:Bucket                   tatumsecurity-bedrock-input-s3                          deleted (0.79s)
 -   └─ aws:cloudwatch:EventRule        TatumInputBucketRule                                    deleted (0.89s)

Resources:
    - 18 deleted

Duration: 1m5s

The resources in the stack have been deleted, but the history and configuration associated with the stack are still maintained.
If you want to remove the stack completely, run `pulumi stack rm pulumi-cloud-ai-exam`.
```

## 스택 삭제

스택까지 완벽하게 삭제하고싶을 경우 실행합니다.

```
$ pulumi stack rm
This will permanently remove the 'pulumi-cloud-ai-exam' stack!
Please confirm that this is what you'd like to do by typing `pulumi-cloud-ai-exam`: pulumi-cloud-ai-exam
Stack 'pulumi-cloud-ai-exam' has been removed!
```
