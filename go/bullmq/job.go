package bullmq

type BackoffStrategy interface {
	Delay(attempts int) (msec int)
}

type BackoffFunc func(attempts int) (msec int)

func (b BackoffFunc) Delay(attempts int) (msec int) {
	return b(attempts)
}

// RemovePolicy represents TS type boolean | number | KeepJobs;
// if Always == true, then always remove, regardless of Count and AgeSeconds
// if Count != 0, keep up to N, for up to AgeSeconds if set
// if AgeSeconds is set, keep up to AgeSeconds. If Count==0, keep infinite, otherwise defer to KeepPolicy
type RemovePolicy struct {
	Always     bool `json:"always"`
	Count      int  `json:"count"`
	AgeSeconds int  `json:"age"`
}

type DefaultJobOptions struct {
	Timestamp        *int            `json:"timestamp,omitempty"`
	Priority         int             `json:"priority,omitempty"`
	Delay            int             `json:"delay,omitempty"`
	Attempts         int             `json:"attempts,omitempty"`
	Backoff          BackoffStrategy `json:"backoff,omitempty"` //TODO: support a custom unmarshaler
	Lifo             bool            `json:"lifo,omitempty"`
	RemoveOnComplete *RemovePolicy   `json:"removeOnComplete,omitempty"`
	RemoveOnFail     *RemovePolicy   `json:"removeOnFail,omitempty"`
	KeepLogs         int             `json:"keepLogs,omitempty"`
	StackTraceLimit  int             `json:"stackTraceLimit,omitempty"`
	SizeLimit        int             `json:"sizeLimit,omitempty"`
}

type BaseJobOptions struct {
	DefaultJobOptions

	// TODO: Repeat
	//Repeat RepeatOptions `json:"repeat,omitempty"`
	//RepeatJobKey string `json:"repeatJobKey,omitempty"`

	JobId  *string `json:"jobId,omitempty"`
	Parent *JobRef `json:"parent,omitempty"`

	/**
	 * Internal property used by repeatable jobs.
	 */
	PrevMillis *int `json:"prevMillis,omitempty"`
}

type JobRef struct {
	Id string `json:"id"`

	/**
	 * It includes the prefix, the namespace separator :, and queue name.
	 * @see https://www.gnu.org/software/gawk/manual/html_node/Qualified-Names.html
	 */
	Queue string `json:"queue"`
}
