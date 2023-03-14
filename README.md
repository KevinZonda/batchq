# BatchQ

BatchQ provides a simple way to aggregate multiple jobs into a single batch job.
This is useful for running jobs on a cluster where you have a limited number of jobs that can be run at once.

i.e. for some APIs, its QPS is 10rps, but it allows you to send a request with
several jobs. You have 1000 jobs to run, you can use BatchQ to aggregate 10
jobs into a single job, and run it every 0.1 seconds.