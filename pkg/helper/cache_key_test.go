package helper

import "testing"

func TestCacheKeyLockTrxCreditUser(t *testing.T) {
	tests := []struct {
		userID   string
		expected string
	}{
		{"user123", "xyz-skilltest:trx-credit-user:user123"},
		{"abc456", "xyz-skilltest:trx-credit-user:abc456"},
		{"testUser", "xyz-skilltest:trx-credit-user:testUser"},
		{"", "xyz-skilltest:trx-credit-user:"},
	}

	for _, tt := range tests {
		t.Run("CacheKeyLockTrxCreditUser", func(t *testing.T) {
			result := CacheKeyLockTrxCreditUser(tt.userID)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for userID %v", tt.expected, result, tt.userID)
			}
		})
	}
}
