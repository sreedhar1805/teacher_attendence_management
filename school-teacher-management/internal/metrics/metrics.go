package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// =========================
	// HTTP METRICS (EXISTING)
	// =========================

	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// =========================
	// TEACHER METRICS
	// =========================

	TeachersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "teachers_created_total",
			Help: "Total number of teachers created",
		},
	)

	TeachersTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "teachers_total",
			Help: "Current total number of teachers",
		},
	)

	// =========================
	// ATTENDANCE METRICS
	// =========================

	AttendanceCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "attendance_created_total",
			Help: "Total number of attendance records created",
		},
	)

	AttendanceCheckInTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "attendance_checkin_total",
			Help: "Total number of check-ins",
		},
	)

	AttendanceCheckOutTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "attendance_checkout_total",
			Help: "Total number of check-outs",
		},
	)

	AttendanceTodayCheckedIn = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "attendance_today_checked_in",
			Help: "Number of teachers checked in today",
		},
	)
)

// Register all metrics here
func Register() {
	prometheus.MustRegister(
		// HTTP
		HttpRequestsTotal,
		HttpRequestDuration,

		// Teachers
		TeachersCreatedTotal,
		TeachersTotal,

		// Attendance
		AttendanceCreatedTotal,
		AttendanceCheckInTotal,
		AttendanceCheckOutTotal,
		AttendanceTodayCheckedIn,
	)
}
