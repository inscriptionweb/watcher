package sender

import "github.com/pkg/sftp"
import "golang.org/x/crypto/ssh"
import "io/ioutil"
import "log"
import "strings"

type Sender struct {
	sftpClient *sftp.Client
	remotePath string
	localPath  string
}

// Constructor
func NewSender(username string, ip string, keyFile string, localPath string, remotePath string) *Sender {
	privateKeyDatas, _ := ioutil.ReadFile(keyFile)

	signer, error := ssh.ParsePrivateKey([]byte(privateKeyDatas))

	if error != nil {
		log.Fatalf("Failed to parse private key: " + error.Error())
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	var sftpClient *sftp.Client

	if client, err := ssh.Dial("tcp", ip+":22", config); err == nil {
		if sftpClientTmp, err := sftp.NewClient(client); err == nil {
			sftpClient = sftpClientTmp
		} else {
			log.Fatalf("Failed to connect to remote host: " + err.Error())
		}
	}

	return &Sender{
		sftpClient: sftpClient,
		localPath:  localPath,
		remotePath: remotePath,
	}

}

// Send a file remotely
func (s *Sender) Send(files *[]string) {

	for _, filename := range *files {
		cleanedLocalPath := strings.TrimRight(s.localPath, "/")
		cleanedRemotePath := strings.TrimRight(s.remotePath, "/")
		fileWithoutSuffix := strings.TrimPrefix(filename, cleanedLocalPath+"/")
		remoteFileName := cleanedRemotePath + "/" + strings.TrimPrefix(filename, cleanedLocalPath+"/")
		localFileDatas, _ := ioutil.ReadFile(filename)

		chunks := strings.Split(fileWithoutSuffix, "/")
		chunks = chunks[:len(chunks)-1]

		folders := cleanedRemotePath

		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("Check remote folder exists and connection to remote host is possible")
			}
		}()

		_, folderError := s.sftpClient.Lstat(folders + "/" + strings.Join(chunks, "/"))

		error, ok := folderError.(*sftp.StatusError)

		if ok && error != nil && error.Code == 2 {

			for _, chunk := range chunks {

				folders = folders + "/" + chunk

				log.Print("Create folder " + folders)

				s.sftpClient.Mkdir(folders)
			}
		}

		remoteFile, _ := s.sftpClient.Create(remoteFileName)
		defer remoteFile.Close()
		remoteFile.Write(localFileDatas)

		log.Print("Copy file " + remoteFileName)
	}
}
