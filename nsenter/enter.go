package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
void __attribute__((constructor))  enter_namespace(void) {
	char *wdocker_pid;
	wdocker_pid = getenv("wdocker_pid");
	if (wdocker_pid) {
		//fprintf(stdout, "got wdocker_pid=%s\n", wdocker_pid);
	} else {
		// fprintf(stdout, "missing wdocker_pid env skip nsenter");
		return;
	}
	char *wdocker_cmd;
	wdocker_cmd = getenv("wdocker_cmd");
	if (wdocker_cmd) {
		//fprintf(stdout, "got wdocker_cmd=%s\n", wdocker_cmd);
	} else {
		//fprintf(stdout, "missing wdocker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };
	for (i=0; i<5; i++) {
		sprintf(nspath, "/proc/%s/ns/%s", wdocker_pid, namespaces[i]);
		int fd = open(nspath, O_RDONLY);
		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		} else {
			//fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		}
		close(fd);
	}
	int res = system(wdocker_cmd);
	exit(0);
	return;
}
*/
import "C"
