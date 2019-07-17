// Copyright 2019 RedHat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package localmetrics

import (
	"math"

	awsv1alpha1 "github.com/openshift/aws-account-operator/pkg/apis/aws/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MetricTotalAWSAccounts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_aws_accounts",
		Help: "Report how many accounts have been created in AWS org",
	}, []string{"name"})
	MetricTotalAccountCRs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_account_crs",
		Help: "Report how many account CRs have been created",
	}, []string{"name"})
	MetricTotalAccountCRsUnclaimed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_accounts_crs_unclaimed",
		Help: "Report how many account CRs are unclaimed",
	}, []string{"name"})
	MetricTotalAccountCRsClaimed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_accounts_crs_claimed",
		Help: "Report how many account CRs are claimed",
	}, []string{"name"})
	MetricTotalAccountCRsFailed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_accounts_crs_failed",
		Help: "Report how many account  CRs are failed",
	}, []string{"name"})
	MetricTotalAccountClaimCRs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_aws_account_claim_crs",
		Help: "Report how many account claim CRs have been created",
	}, []string{"name"})
	MetricPoolSizeVsUnclaimed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_pool_size_vs_unclaimed",
		Help: "Report the difference between the pool size and the number of unclaimed account CRs",
	}, []string{"name"})
	MetricTotalAccountPendingVerification = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "aws_account_operator_total_aws_account_pending_verification",
		Help: "Report the number of accounts waiting for enterprise support and EC2 limit increases in AWS",
	}, []string{"name"})

	MetricsList = []prometheus.Collector{
		MetricTotalAWSAccounts,
		MetricTotalAccountCRs,
		MetricTotalAccountCRsUnclaimed,
		MetricTotalAccountCRsClaimed,
		MetricTotalAccountCRsFailed,
		MetricTotalAccountClaimCRs,
		MetricPoolSizeVsUnclaimed,
		MetricTotalAccountPendingVerification,
	}
)

// UpdateAWSMetrics updates all AWS related metrics
func UpdateAWSMetrics(totalAccounts int) {
	MetricTotalAWSAccounts.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(totalAccounts))
}

// UpdateAccountCRMetrics updates all metrics related to Account CRs
func UpdateAccountCRMetrics(accountList *awsv1alpha1.AccountList) {

	unclaimedAccountCount := 0
	claimedAccountCount := 0
	failedAccountCount := 0
	pendingVerificationAccountCount := 0
	for _, account := range accountList.Items {
		if account.Status.Claimed == false {
			if account.Status.State != "Failed" {
				unclaimedAccountCount++
			}
		} else {
			claimedAccountCount++
		}
		if account.Status.State == "Failed" {
			failedAccountCount++
		} else if account.Status.State == "PendingVerification" {
			pendingVerificationAccountCount++
		}
	}

	MetricTotalAccountCRs.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(len(accountList.Items)))
	MetricTotalAccountCRsUnclaimed.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(unclaimedAccountCount))
	MetricTotalAccountCRsClaimed.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(claimedAccountCount))
	MetricTotalAccountPendingVerification.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(pendingVerificationAccountCount))
	MetricTotalAccountCRsFailed.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(failedAccountCount))
}

// UpdateAccountClaimMetrics updates all metrics related to AccountClaim CRs
func UpdateAccountClaimMetrics(accountClaimList *awsv1alpha1.AccountClaimList) {

	MetricTotalAccountClaimCRs.With(prometheus.Labels{"name": "aws-account-operator"}).Set(float64(len(accountClaimList.Items)))
}

// UpdatePoolSizeVsUnclaimed updates the metric that measures the difference between Poolsize and Unclaimed Account CRs
func UpdatePoolSizeVsUnclaimed(poolSize int, unclaimedAccountCount int) {

	metric := math.Abs(float64(poolSize) - float64(unclaimedAccountCount))

	MetricPoolSizeVsUnclaimed.With(prometheus.Labels{"name": "aws-account-operator"}).Set(metric)
}