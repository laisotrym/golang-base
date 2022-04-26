package status_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    
    promo_status "safeweb.app/service/safeweb_lib/status"
)

func TestTemplate_String(t *testing.T) {
    t.Run("200", func(t *testing.T) {
        // Arrange
        expected := promo_status.NewErrorTrace(promo_status.Success, "Khởi tạo giao dịch thành công")
        
        // Act
        actual := promo_status.DefaultSuccessTemplate.With("Khởi tạo giao dịch")
        
        // Assert
        assert.Equal(t, expected, actual)
    })
    t.Run("410", func(t *testing.T) {
        // Arrange
        expected1 := promo_status.NewErrorTrace(promo_status.Unauthorized, "Không thể xác thực chữ ký hoặc thông tin")
        expected2 := promo_status.NewErrorTrace(promo_status.Unauthorized, "IP không được phép truy cập")
        expected3 := promo_status.NewErrorTrace(promo_status.Unauthorized, "Sai checksum")
        
        // Act
        actual1 := promo_status.DefaultUnauthorizedTemplate.With()
        actual2 := promo_status.IPNotAllowedAccessTemplate.With()
        actual3 := promo_status.WrongChecksumTemplate.With()
        
        // Assert
        assert.Equal(t, expected1, actual1)
        assert.Equal(t, expected2, actual2)
        assert.Equal(t, expected3, actual3)
    })
}
