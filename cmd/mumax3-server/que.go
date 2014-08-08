package main

import "log"

func (n *Node) GiveJob(nodeAddr string) string {
	n.lock()
	defer n.unlock()

	jobs := n.jobs
	if len(jobs) == 0 {
		return "" // indicates no job available
	}
	job := jobs[len(jobs)-1]
	jobs = jobs[:len(jobs)-1]

	job.Node = nodeAddr
	n.running = append(n.running, job)
	url := "http://" + n.inf.Addr + "/fs/" + job.File
	log.Println("give job", url, "->", nodeAddr)
	return url
}

func (n *Node) AddJob(fname string) {
	n.lock()
	defer n.unlock()
	log.Println("Push job:", fname)
	n.jobs = append(n.jobs, Job{File: fname})
}

// convert []*Job to []string (job names),
// for http status etc.
func copyJobs(jobs []Job) []Job {
	return append([]Job{}, jobs...)
}
