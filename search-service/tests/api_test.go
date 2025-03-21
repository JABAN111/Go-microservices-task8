package words_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const address = "http://localhost:28080"

var client = http.Client{
	Timeout: 5 * time.Minute,
}

func TestPreflight(t *testing.T) {
	require.Equal(t, true, true)
}

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

func TestPing(t *testing.T) {
	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err, "cannot ping")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "wrong status")

	var reply PingResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&reply))
	require.Equal(t, "ok", reply.Replies["words"], "no words running")
	require.Equal(t, "ok", reply.Replies["update"], "no db running")
}

type UpdateStats struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

type UpdateStatus struct {
	Status string `json:"status"`
}

func prepare(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, address+"/api/db", nil)
	require.NoError(t, err, "cannot make request")
	resp, err := client.Do(req)
	require.NoError(t, err, "could not send clean up command")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	st := stats(t)
	require.Equal(t, 0, st.ComicsFetched)
	require.True(t, st.ComicsTotal > 3000, "there are more than 3000 comics in XKCD")
	require.Equal(t, 0, st.WordsTotal)
	require.Equal(t, 0, st.WordsUnique)

	require.Equal(t, "idle", status(t))
}

func TestEmptyDB(t *testing.T) {
	prepare(t)
}

func TestUpdate(t *testing.T) {
	prepare(t)
	var wg sync.WaitGroup
	wg.Add(3)
	var res1, res2 int
	var res3 string
	go func() {
		res1 = update(t)
		wg.Done()
	}()
	go func() {
		res2 = update(t)
		wg.Done()
	}()
	go func() {
		time.Sleep(1 * time.Second)
		res3 = status(t)
		wg.Done()
	}()
	wg.Wait()

	fmt.Println(res1, res2)
	require.True(t,
		res1 == http.StatusOK && res2 == http.StatusAccepted ||
			res2 == http.StatusOK || res1 == http.StatusAccepted,
		"wrong statuses from concurrent updates, expect ok && accepted",
	)
	// FIXME: (если все сервисы подключены через контейнеры кроме api, то все тесты проходят)
	// Если api запущен именно как контейнер, то eof ошибка
	// возможно, но у меня не получилось  https://stackoverflow.com/questions/76848446/running-docker-containers-on-mac-apple-silicon-m2-platform-compatibility-and-ex

	require.Equal(t, "running", res3, "need running status while update")
	st := stats(t)
	require.Equal(t, st.ComicsTotal, st.ComicsFetched)
	require.True(t, st.ComicsTotal > 3000, "there are more than 3000 comics in XKCD")
	require.True(t, 1000 < st.WordsTotal, "not enough total words in DB")
	require.True(t, 100 < st.WordsUnique, "not enough unique words in DB")

	prepare(t)
}

// кидает ошибку и зависает на ней(то есть тест не завершается аварийно)
// opt/homebrew/opt/go/libexec/src/runtime/asm_arm64.s:1223

// Фиксится странным костылем, а именно запуском `update` через горутину, то есть отдавать ответы в асинхронном режиме
// по логам видно, что http отправляет свой ответ, но почему-то ответ, который ждет конца update разваливается
func update(t *testing.T) int {
	req, err := http.NewRequest(http.MethodPost, address+"/api/db/update", nil)
	require.NoError(t, err, "cannot make request")
	resp, err := client.Do(req)
	fmt.Println("resp there", resp)
	// resp there &{202 Accepted 202 HTTP/1.1 1 1 map[Content-Length:[0] Date:[Thu, 20 Mar 2025 04:00:37 GMT]] 0x140001a8240 0 [] false false map[] 0x140001aea00 <nil>}
	//resp there <nil>
	if err != nil {
		panic(err)
	}
	require.NoError(t, err, "could not send update command")
	defer resp.Body.Close()
	return resp.StatusCode
}

func status(t *testing.T) string {
	resp, err := client.Get(address + "/api/db/status")
	require.NoError(t, err, "could not get status")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var status UpdateStatus
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&status), "cannot decode")
	return status.Status
}

func stats(t *testing.T) UpdateStats {
	resp, err := client.Get(address + "/api/db/stats")
	require.NoError(t, err, "could not get stats")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var stats UpdateStats
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&stats), "cannot decode")
	return stats
}
