package fileobjects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	minio "github.com/minio/minio-go/v7"

	"github.com/fils/goobjectweb/internal/fileactions"
)

// FileObjects pulls the objects from the object store
func FileObjects(mc *minio.Client, bucket, prefix, domain string, w http.ResponseWriter, r *http.Request) {
	var object string
	key := fmt.Sprintf("%s", r.URL.Path)

	// TODO review this hack...
	if key == "" {
		key = "index.html"
	}
	if strings.HasSuffix(key, "/") {
		key = key + "index.html"
	}

	// TODO if key has no extension, make it text/html
	m := fileactions.MimeByType(filepath.Ext(key))
	// TODO remove hack  (need to check if no suffic, make text/html for react)
	m = "text/html"
	w.Header().Set("Content-Type", m)
	log.Printf("%s: %s \n", key, m)

	// TODO ocd update March 21  CSDCO
	// caution, is this needed ??????
	fmt.Println(prefix)
	//newprefix := strings.Replace(prefix, "website", "csdco", 1)
	//fmt.Printf("New prefix: %s  old prefix: %s", newprefix, prefix)
	//prefix = newprefix

	object = fmt.Sprintf("%s/%s", prefix, key)
	//object = fmt.Sprintf("%s/website/%s", prefix, key)
	//object = fmt.Sprintf("website/%s", key)

	log.Printf("bucket: %s   object: %s", bucket, object)

	// check our object is there first....
	oi, err := mc.StatObject(context.Background(), bucket, object, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Error: %v Size: %d", err, oi.Size)
		//w.WriteHeader(http.StatusNotFound)
		//http.Redirect(w, r, "/", 303) // HACK looking at react -s no -s behavior
		// March 2022, react hack back to index.html
		object = fmt.Sprintf("%s/%s", prefix, "index.html")
	}

	fo, err := mc.GetObject(context.Background(), bucket, object, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error getobjt: %s  \n err: %s", object, err)
		return
	}

	size, err := io.Copy(w, fo) // todo need to stream write from s3 reader...  not copy
	if err != nil {
		log.Println("Issue with writing file to http response")
		log.Println(err)
	}

	log.Print(size)
}
