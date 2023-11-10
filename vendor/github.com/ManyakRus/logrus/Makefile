SERVICEURL=github.com/sirupsen/logrus

NEW_REPO=github.com/ManyakRus/logrus

newrepo:
	find -name *.go -not -path "*/vendor/*"|xargs sed -i 's+$(SERVICEURL)+$(NEW_REPO)+g'
