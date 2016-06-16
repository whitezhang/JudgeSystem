#include <sys/resource.h>
#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/user.h>
#include <sys/syscall.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>

const int LANGC = 2;
const int LANGJAVA = 3;
const int DEBUG = 1;
const int STD_MB = 1024*1024;
//const int STD_MB = 10;
const int BUFFER_SIZE = 1024;

void run_solution(int runid, int clientid) {
    char buf[BUFFER_SIZE], runidstr[BUFFER_SIZE];
    struct rlimit LIM;
    ptrace(PTRACE_TRACEME, 0, NULL, NULL);
    // max time of cpu
    LIM.rlim_max = 800;
    LIM.rlim_cur = 800;
    setrlimit(RLIMIT_CPU, &LIM);

    // max length of file can be built
    LIM.rlim_max = 1 * STD_MB;
    LIM.rlim_cur = 1 * STD_MB;
    setrlimit(RLIMIT_FSIZE, &LIM);

    // max virtual memory, heap
    LIM.rlim_max = STD_MB * 10;
    LIM.rlim_cur = STD_MB * 10;
    //LIM.rlim_max = 100;
    //LIM.rlim_cur = 100;
    setrlimit(RLIMIT_AS, &LIM);

    // max progress for user
    LIM.rlim_max = 200;
    LIM.rlim_cur = 200;
    setrlimit(RLIMIT_NPROC, &LIM);

    LIM.rlim_cur = STD_MB * 20;
    LIM.rlim_cur = STD_MB * 20;
    setrlimit(RLIMIT_STACK, &LIM);

    LIM.rlim_cur = STD_MB * 10;
    LIM.rlim_cur = STD_MB * 10;
    setrlimit(RLIMIT_DATA, &LIM);

    sprintf(runidstr, "%d", runid);
    sprintf(buf, "%d", clientid);
    //*/

    // run client
    //execl("./test", "test", "" , NULL);
    //execl("./main < in.txt", "main", "" , NULL);
    execl("./main", "main", "" , NULL);
}

int get_proc_status(int pid, const char * mark) {
    FILE * pf;
    char fn[BUFFER_SIZE], buf[BUFFER_SIZE];
    int ret = 0;
    sprintf(fn, "/proc/%d/status", pid);
    pf = fopen(fn, "re");
    int m = strlen(mark);
    while (pf && fgets(buf, BUFFER_SIZE - 1, pf)) {

        buf[strlen(buf) - 1] = 0;
        if (strncmp(buf, mark, m) == 0) {
            sscanf(buf + m + 1, "%d", &ret);
        }
    }
    if (pf)
        fclose(pf);
    return ret;
}

int get_page_fault_mem(struct rusage & ruse, pid_t & pidApp) {
    int m_vmpeak, m_vmdata, m_minflt;
    m_minflt = ruse.ru_minflt * getpagesize();
    if(DEBUG) {
        // Memory of this process
        m_vmpeak = get_proc_status(pidApp, "VmPeak:");
        // Memory of data segment
        m_vmdata = get_proc_status(pidApp, "VmData:");
        printf("VmPeak:%d KB VmData:%d KB minflt:%d KB\n", m_vmpeak, m_vmdata, m_minflt >> 10);
    }
    return m_minflt;
}

int watch_solution(pid_t pidApp, int lang, int topmemory, int mem_lmt, int time_lmt) {
    int status, exitcode, sig;
    struct rusage ruse;
    struct user_regs_struct reg;
    int tempmemory, usedtime;
    usedtime = 0;

    if(topmemory == 0) {
        // Memory being used by this process
        topmemory = get_proc_status(pidApp, "VmRss:") << 10;
        if(DEBUG) {
            printf("VMRss:%d KB\n", topmemory);
        }
    }

    while(1) {
        wait4(pidApp, &status, 0, &ruse);
        if (lang == LANGJAVA) {
            tempmemory = get_page_fault_mem(ruse, pidApp);
        }
        else {
            tempmemory = get_proc_status(pidApp, "VmPeak:") << 10;
            if(DEBUG) {
                printf("VmPeak:%d KB\n", tempmemory);
            }
        }

        if(tempmemory > topmemory) {
            topmemory = tempmemory;
        }
        printf("%.2fMB\t%dMB\n", topmemory*1.0/STD_MB, mem_lmt);
        if(topmemory > mem_lmt * STD_MB) {
            if(DEBUG) {
                printf("Out of memory %d\n", topmemory);
            }
            ptrace(PTRACE_KILL, pidApp, NULL, NULL);
            break;
        }

        // 0 if the subprocess successed
        if(WIFEXITED(status)) {
            break;
        }
        exitcode = WEXITSTATUS(status);
        // if the subprocess terminted by signal
        if(WIFSIGNALED(status)) {
            sig = WTERMSIG(status);
            if(DEBUG) {
                printf("WTERMSIG=%d\n", sig);
                psignal(sig, NULL);
            }
            break;
        }
        usedtime += (ruse.ru_utime.tv_sec * 1000 + ruse.ru_utime.tv_usec / 1000);
        usedtime += (ruse.ru_stime.tv_sec * 1000 + ruse.ru_stime.tv_usec / 1000);
        if(usedtime >= time_lmt) {
            ptrace(PTRACE_KILL, pidApp, NULL, NULL);
            break;
        }

        ptrace(PTRACE_GETREGS, pidApp, NULL, &reg);
        //printf("%ld\n", (long)reg.orig_rax);
        ptrace(PTRACE_SYSCALL, pidApp, NULL, NULL);
    }
    return 0;
}

int main() {
    int mem_lmt = 32;   // 32MB
    int time_lmt = 1000;    // 1000ms
    int topmemory = 0;

    pid_t pidApp = fork();
    if(pidApp == 0) {
        run_solution(1, 1);
    }
    else {
        watch_solution(pidApp, LANGC, topmemory, mem_lmt, time_lmt);
    }
}
