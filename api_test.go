package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockJobClient struct {
	hadoopJobClient
	mock.Mock
}

type mockPersistedJobClient struct {
	s3JobClient
	mock.Mock
}

func (m *mockPersistedJobClient) FetchJob(id string) (*job, error) {
	args := m.Called(id)
	detail := args.Get(0)
	if detail != nil {
		return args.Get(0).(*job), args.Error(1)
	}
	return nil, args.Error(1)
}

/**
 * sets a jobTracker to main.jts
 */
func setJobTracker(client RecentJobClient) *jobTracker {
	jts = make(map[string]*jobTracker)
	var jt = newJobTracker("foo", client)
	jts["testCluster"] = jt
	return jt
}

func TestGetNonExistentJob(t *testing.T) {
	id := "nonexistentjob"
	mockStorageClient := new(mockPersistedJobClient)
	mockStorageClient.On("FetchJob", id).Return(nil, fmt.Errorf("Bad"))
	persistedJobClient = mockStorageClient

	res := getJob(id)

	assert.Nil(t, res)
}

func TestGetJobFromMemory(t *testing.T) {
	var id = "job_from_memory"
	var job = &job{}
	var jt = setJobTracker(new(mockJobClient))
	jt.jobs[jobID(id)] = job

	res := getJob(id)

	assert.Equal(t, res, job)
}

func TestGetJobFromS3(t *testing.T) {
	var id = "job_from_s3"
	var myjob = &job{}
	mockStorageClient := new(mockPersistedJobClient)
	mockStorageClient.On("FetchJob", id).Return(myjob, nil)
	persistedJobClient = mockStorageClient

	res := getJob(id)

	assert.Equal(t, myjob, res)
}
