package sender

import "github.com/pkg/sftp"
import "golang.org/x/crypto/ssh"
import "io/ioutil"
import "github.com/Sirupsen/logrus"
import "strings"
import "fmt"

const (
	NO_ERROR                                                       = 0
	FAILED_TO_PARSE_PRIVATE_KEY                                    = 1
	FAILED_TO_CONNECT_TO_REMOTE_HOST                               = 2
	FAILED_TO_CONNECT_TO_REMOTE_HOST_OR_REMOTE_FOLDER_DOESNT_EXIST = 3
)

type SenderError struct {
	code uint32
}

func (e SenderError) Error() string {
	return fmt.Sprintf("An error occured : %v", e.CodeString())
}

func (e SenderError) CodeString() string {
	switch e.code {
	case FAILED_TO_PARSE_PRIVATE_KEY:
		return "FAILED_TO_PARSE_PRIVATE_KEY"
	case FAILED_TO_CONNECT_TO_REMOTE_HOST:
		return "FAILED_TO_CONNECT_TO_REMOTE_HOST"
	case FAILED_TO_CONNECT_TO_REMOTE_HOST_OR_REMOTE_FOLDER_DOESNT_EXIST:
		return "FAILED_TO_CONNECT_TO_REMOTE_HOST_OR_REMOTE_FOLDER_DOESNT_EXIST"
	}

	return "NO_ERROR"
}

func (e SenderError) CodeInteger() uint32 {
	return e.code
}

type Sender struct {
	sftpClient *sftp.Client
	remotePath string
	localPath  string
	logger     *logrus.Logger
}

// Constructor
func NewSender(username string, ip string, keyFile string, localPath string, remotePath string, logger *logrus.Logger) (*Sender, SenderError) {
	privateKeyDatas, _ := ioutil.ReadFile(keyFile)

	signer, error := ssh.ParsePrivateKey([]byte(privateKeyDatas))

	if error != nil {
		senderError := SenderError{
			code: FAILED_TO_PARSE_PRIVATE_KEY,
		}

		logger.WithFields(logrus.Fields{
			"code": senderError.CodeInteger(),
		}).Error(senderError.CodeString())

		return &Sender{}, senderError
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
			senderError := SenderError{
				code: FAILED_TO_CONNECT_TO_REMOTE_HOST,
			}

			logger.WithFields(logrus.Fields{
				"code": senderError.CodeInteger(),
			}).Error(senderError.CodeString())

			return &Sender{}, senderError
		}
	}

	return &Sender{
		sftpClient: sftpClient,
		localPath:  localPath,
		remotePath: remotePath,
		logger:     logger,
	}, SenderError{}

}

// Send a file remotely
func (s *Sender) Send(files *[]string) SenderError {

	senderError := SenderError{
		code: NO_ERROR,
	}

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
				senderError = SenderError{
					code: FAILED_TO_CONNECT_TO_REMOTE_HOST_OR_REMOTE_FOLDER_DOESNT_EXIST,
				}

				s.logger.WithFields(logrus.Fields{
					"code": senderError.CodeInteger(),
				}).Error(senderError.CodeString())
			}
		}()

		if senderError.CodeInteger() != NO_ERROR {
			return senderError
		}

		_, folderError := s.sftpClient.Lstat(folders + "/" + strings.Join(chunks, "/"))

		if error, ok := folderError.(*sftp.StatusError); ok && error != nil && error.Code == 2 {

			for _, chunk := range chunks {

				folders = folders + "/" + chunk

				s.sftpClient.Mkdir(folders)

				s.logger.WithFields(logrus.Fields{
					"folders": folders,
				}).Info("Folder created")
			}
		}

		remoteFile, _ := s.sftpClient.Create(remoteFileName)
		defer remoteFile.Close()
		remoteFile.Write(localFileDatas)

		s.logger.WithFields(logrus.Fields{
			"file": remoteFileName,
		}).Info("File copied")
	}

	return senderError
}
