package safeweb_lib_helper

func MinFloat(first float64, rest ...float64) float64 {
    ans := first
    for _, item := range rest {
        if item < ans {
            ans = item
        }
    }
    return ans
}

func MaxFloat(first float64, rest ...float64) float64 {
    ans := first
    for _, item := range rest {
        if item > ans {
            ans = item
        }
    }
    return ans
}

func MinInt(first int, rest ...int) int {
    ans := first
    for _, item := range rest {
        if item < ans {
            ans = item
        }
    }
    return ans
}
