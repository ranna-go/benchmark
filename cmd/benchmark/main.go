package main

import (
	"flag"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/ranna-go/benchmark/pkg/workerpool"
	ranna "github.com/ranna-go/ranna/pkg/client"
	models "github.com/ranna-go/ranna/pkg/models"
	"github.com/sirupsen/logrus"
)

var (
	endpoint    = flag.String("e", "", "ranna API endpoint")
	loglevel    = flag.Int("loglevel", int(logrus.InfoLevel), "logger level")
	number      = flag.Int("n", 5, "number of total requests")
	parallel    = flag.Int("p", 5, "parallel running requests")
	snippets    = flag.String("snippets", "./snippets", "snippets directory")
	prettyPrint = flag.Bool("pp", false, "pretty print results")
)

type timing struct {
	ReqDur  time.Duration
	ExecDur time.Duration
}

func must(err error) {
	if err != nil {
		logrus.WithError(err).Fatal()
	}
}

func main() {
	flag.Parse()
	logrus.SetLevel(logrus.Level(*loglevel))
	logrus.SetFormatter(&logrus.TextFormatter{})

	logrus.Info("read source snippets ...")
	sources, err := readSourceFiles(*snippets)
	must(err)

	if *endpoint == "" {
		logrus.Fatal("endpoint must be specified")
	}

	client, err := ranna.New(ranna.Options{
		Endpoint: *endpoint,
	})
	must(err)

	logrus.WithField("parallel", *parallel).WithField("n", *number).Info("initializing worker pool ...")
	wp := workerpool.New(*parallel)

	timings := make([]*timing, 0)
	go func() {
		for {
			v := <-wp.Results()
			if t, ok := v.(*timing); ok {
				timings = append(timings, t)
			}
		}
	}()
	for i := 0; i < *number; i++ {
		wp.Push(job(client, pickReq(sources)))
	}
	wp.Close()
	wp.WaitBlocking()

	logrus.Info("sleep for 1s to finish collection")
	time.Sleep(1 * time.Second)

	av, n := calcAverageTimings(timings)
	logrus.WithFields(logrus.Fields{
		"av_req":     av.ReqDur,
		"av_exec":    av.ExecDur,
		"successful": n,
		"errros":     *number - n,
	}).Info("bench result")

	if *prettyPrint {
		tab := table.NewWriter()
		tab.SetOutputMirror(os.Stdout)
		tab.SetStyle(table.StyleLight)
		tab.AppendHeader(table.Row{"Used Benchmark Parameters"})
		tab.AppendSeparator()
		tab.AppendRows([]table.Row{
			{"Endpoint", *endpoint},
			{"# Snippets", len(sources)},
			{"# Requests", *number},
			{"# Parallel Reqeusts", *parallel},
		})
		tab.Render()

		tab = table.NewWriter()
		tab.SetOutputMirror(os.Stdout)
		tab.SetStyle(table.StyleLight)
		tab.AppendHeader(table.Row{"Benchmark Results"})
		tab.AppendSeparator()
		tab.AppendRows([]table.Row{
			{"# Successful", n},
			{"% Successful", n / *number * 100},
			{"# Erroneous", *number - n},
			{"% Erroneous", (*number - n) / *number * 100},
			{"Average Request to Response Time", av.ReqDur},
			{"Average Execution Time", av.ExecDur},
		})
		tab.Render()
	}
}

func pickReq(arr []*models.ExecutionRequest) *models.ExecutionRequest {
	return arr[rand.Intn(len(arr))]
}

func job(client ranna.Client, req *models.ExecutionRequest) workerpool.Job {
	return func(workerId int, params ...interface{}) interface{} {
		logrus.WithField("worker", workerId).WithField("lang", req.Language).Info("exec")
		t := &timing{}
		now := time.Now()
		res, err := client.Exec(*req)
		t.ReqDur = time.Since(now)
		t.ExecDur = time.Millisecond * time.Duration(res.ExecTimeMS)
		if err != nil {
			t = nil
			logrus.WithField("worker", workerId).WithField("lang", req.Language).WithError(err).Error("exec failed")
		} else {
			logrus.WithField("worker", workerId).WithField("lang", req.Language).Info("exec finished")
		}
		return t
	}
}

func readSourceFiles(snippetsDir string) (res []*models.ExecutionRequest, err error) {
	files, err := os.ReadDir(snippetsDir)
	if err != nil {
		return
	}
	res = make([]*models.ExecutionRequest, len(files))
	i := 0
	var content []byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		lang := path.Base(f.Name())[0 : len(f.Name())-len(path.Ext(f.Name()))]
		content, err = os.ReadFile(path.Join(snippetsDir, f.Name()))
		if err != nil {
			return
		}
		res[i] = &models.ExecutionRequest{
			Language: lang,
			Code:     string(content),
		}
		i++
	}
	res = res[:i]
	return
}

func calcAverageTimings(arr []*timing) (timing, int) {
	res := timing{}
	i := 0
	for _, t := range arr {
		if t != nil {
			res.ReqDur += t.ReqDur
			res.ExecDur += t.ExecDur
			i++
		}
	}
	if i > 0 {
		res.ReqDur /= time.Duration(i)
		res.ExecDur /= time.Duration(i)
	}
	return res, i
}
