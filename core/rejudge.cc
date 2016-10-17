#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <dirent.h>
#include <unistd.h>
#include <time.h>
#include <stdarg.h>
#include <ctype.h>
#include <sys/wait.h>
#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/user.h>
#include <sys/syscall.h>
#include <sys/time.h>
#include <sys/resource.h>
#include <sys/signal.h>
#include <sys/stat.h>
#include <unistd.h>
#include <assert.h>
#include <ctime>
#include <fstream>

#include "mongo/client/dbclient.h"
#include "mongo/client/dbclientinterface.h"
#include <sys/stat.h>

using namespace std;

const int BATCH_SIZE = 128;

const int LANGC = 1;
const int LANGCPP = 2;
const int LANGJAVA = 3;
const int DEBUG = 1;
const int STD_MB = 1024*1024;
const int BUFFER_SIZE = 1024;

const int SUCCESS = 0;
const int ERROR = 1;

const int GLOBAL_ID = 1536;
const char *RM_PATH = "/home/judger/worksystem/runtime";
const char *DB_HOST = "127.0.0.1";
const char *DATA_PATH = "/home/judger/worksystem/data";

const int F_AC = 0;
const int F_WA = 1;
const int F_RE = 2;
const int F_TLE = 3;
const int F_PE = 4;
const int F_CE = 5;

const char *F_SET[6] = {"Accepted", "Wrong Answer", "Runtime Error", "Time Limited Error", "Presentation Error", "Compile Error"};

typedef struct problem_runtime {
    char *pathname;
    char *filename;
    char *lang;
    char *oid;
    int pid;
}problem_runtime;

int execute_cmd(const char *fmt, ...) {
    char cmd[BUFFER_SIZE];
    int ret = 0;
    va_list ap;
    va_start(ap, fmt);
    vsprintf(cmd, fmt, ap);
    ret = system(cmd);
    va_end(ap);
    return ret;
}

void run_solution(int runid, int clientid, const char* data_in) {
    nice(19);
    freopen(data_in, "r", stdin);
    freopen("user.out", "w", stdout);
    freopen("err.out", "w", stderr);

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

    // run client
    //execl("./test", "test", "" , NULL);
    execl("./Main", "Main", "" , NULL);
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

int compile(int lang) {
    int pid;
    const char *CP_C[] = {"gcc", "Main.c", "-o", "Main", "-fno-asm", "-Wall",
                        "-lm", "--static", "-std=c99", "-DONLINE_JUDGE", NULL};
    const char *CP_CPP[] = {"g++", "Main.cc", "-o", "Main", "-fno-asm", "-Wall",
                        "-lm", "--static", "-std=c99", "-DONLINE_JUDGE", NULL};
    pid = fork();
    if(pid == 0) {
        struct rlimit LIM;
        LIM.rlim_max = 60;
        LIM.rlim_cur = 60;
        setrlimit(RLIMIT_CPU, &LIM);
        alarm(60);
        LIM.rlim_max = 100 * STD_MB;
        LIM.rlim_cur = 100 * STD_MB;
        setrlimit(RLIMIT_FSIZE, &LIM);

        if(lang == LANGJAVA) {

        }
        else {
            LIM.rlim_max = STD_MB << 10;
            LIM.rlim_cur = STD_MB << 10;
        }
        setrlimit(RLIMIT_AS, &LIM);
        //execute_cmd("chown judge *");
        /*
        while(setgid(GLOBAL_ID) != 0) sleep(1);
        while(setuid(GLOBAL_ID) != 0) sleep(1);
        while(setresuid(GLOBAL_ID, GLOBAL_ID, GLOBAL_ID) != 0) sleep(1);
        */

        switch(lang) {
            case LANGC:
                execvp(CP_C[0], (char * const *) CP_C);
                break;
            case LANGCPP:
                execvp(CP_CPP[0], (char * const *) CP_CPP);
                break;
            default:
                printf("Nothing to do in compiling progress\n");
                break;
        }
        if(DEBUG) {
            printf("Compile finished\n");
        }
        exit(0);
    }
    else {
        int status = 0;
        waitpid(pid, &status, 0);
        if(DEBUG) {
            printf("Compiling status=%d\n", status);
        }
        return status;
    }
}

int watch_runtime(pid_t pidApp, int lang, int topmemory, int mem_lmt, int time_lmt, int *mem_peak, int *time_used) {
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
        if(tempmemory > *mem_peak) {
            *mem_peak = tempmemory;
        }

        if(tempmemory > topmemory) {
            topmemory = tempmemory;
        }
        //printf("%.2fMB\t%dMB\n", topmemory*1.0/STD_MB, mem_lmt);
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
        *time_used = usedtime;
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

// out, user
void find_next_nonspace(int &c1, int &c2, FILE *&f1, FILE *&f2, int &ret) {
    while((isspace(c1) || isspace(c2))) {
        if(c1 != c2) {
            if(c2 == EOF) {
                do {
                    c1 = fgetc(f1);
                }while(isspace(c1));
                continue;
            }
            else if(c1 == EOF) {
                do {
                    c2 = fgetc(f2);
                }while(isspace(c2));
            }
            else if ((c1 == '\r' && c2 == '\n')) {
                c1 = fgetc(f1);
            }
            else {
                ret = F_PE;
            }
        }
        if(isspace(c1)) {
            c1 = fgetc(f1);
        }
        if(isspace(c2)) {
            c2 = fgetc(f2);
        }
    }
}

int compare_solution(const char* file1, const char* file2) {
    int ret = F_AC;
    int c1, c2;
    FILE *f1, *f2;
    f1 = fopen(file1, "re");
    f2 = fopen(file2, "re");
    if(!f1 || !f2) {
        ret = F_RE;
    }
    else {
        for(;;) {
            c1 = fgetc(f1);
            c2 = fgetc(f2);
            find_next_nonspace(c1, c2, f1, f2, ret);
            for(;;) {
                while((!isspace(c1) && c1) || (!isspace(c2) && c2)) {
                    if(c1 == EOF && c2 == EOF) {
                        goto end;
                    }
                    if(c1 == EOF || c2 == EOF) {
                        break;
                    }
                    if(c1 != c2) {
                        ret = F_WA;
                        goto end;
                    }
                    c1 = fgetc(f1);
                    c2 = fgetc(f2);
                }
                find_next_nonspace(c1, c2, f1, f2, ret);
                if(c1 == EOF && c2 == EOF) {
                    goto end;
                }
                if(c1 == EOF || c2 == EOF) {
                    ret = F_WA;
                    goto end;
                }
                if((c1 == '\n' || !c1) && (c2 == '\n' || !c2)) {
                    break;
                }
            }
        }
    }
    end:
    // TODO
    if(f1) fclose(f1);
    if(f2) fclose(f2);
    return ret;
}

int judge_solution(int pid) {
    char std_out_pathname[128];
    sprintf(std_out_pathname, "%s/%d/std.out", DATA_PATH, pid);
    int cmp_ret = compare_solution(std_out_pathname, "user.out");
    return cmp_ret;

    switch(cmp_ret) {
        case F_AC:
            break;
        case F_WA:
            break;
        case F_RE:
            break;
        case F_TLE:
            break;
        case F_PE:
            break;
        default:
            break;
    }
}

void mkdir_if_not_exist(char *pathname) {
    mkdir(pathname, 0777);
}

void mv_code2runtime(problem_runtime *pr_list[BATCH_SIZE]) {
    int index = 0;
    mongo::DBClientConnection c;
    mongo::BSONObj cur_obj;
    mongo::BSONElement cur_ele;
    mongo::OID objectID;
    c.connect(DB_HOST);
    auto_ptr<mongo::DBClientCursor> cursor = c.query("unsys.runtimestatus", mongo::BSONObj());
    while(cursor->more()) {
        time_t t;
        int timestamp = time(&t);
        char pathname[128];
        sprintf(pathname, "%s/%d_%d/", RM_PATH, timestamp, index);

        int pid;
        const char *code, *lang, *oid;

        cur_obj = cursor->next();
        cur_obj.getObjectID(cur_ele);
        objectID = cur_ele.__oid();

        oid = objectID.toString().c_str();
        pid = cur_obj.getIntField("pid");
        code= cur_obj.getStringField("code");
        lang = cur_obj.getStringField("lang");

        mkdir_if_not_exist(pathname);
        strcpy(pr_list[index]->pathname, pathname);
        strcpy(pr_list[index]->lang, lang);
        strcpy(pr_list[index]->oid, oid);
        pr_list[index]->pid = pid;
        if(!strcmp(lang, "c")) {
            strcat(pathname, "Main.c");
        }
        else if(!strcmp(lang, "CPP")) {
            strcat(pathname, "Main.cpp");
        }
        ofstream fout(pathname);
        fout << code;
        fout << flush;
        fout.close();

        strcpy(pr_list[index]->filename, pathname);

        //c.remove("unsys.submitqueue", cur_obj);
        index++;
        if(index > BATCH_SIZE) {
            break;
        }
    }
}

int process_smtqueue(problem_runtime *pr_list[BATCH_SIZE]) {
    mongo::client::initialize();
    try {
        mv_code2runtime(pr_list);
    } catch(const mongo::DBException &e) {
        printf("connection err\n");
    }
    return EXIT_SUCCESS;
}

void update_solution(const char *oid, const int jdg_ret, const int mem_peak, const int time_used) {
    mongo::client::initialize();
    try {
        mongo::DBClientConnection c;

        char c_mem_peak[64], c_time_used[64];
        sprintf(c_mem_peak, "%d", mem_peak);
        sprintf(c_time_used, "%d", time_used);

        c.connect(DB_HOST);
        c.update("unsys.runtimestatus", BSON("_index" << oid), BSON("$set" << BSON( "status"<< F_SET[jdg_ret]<< "memory"<< c_mem_peak<< "time"<< c_time_used)));
    } catch(const mongo::DBException &e) {
        printf("connection err\n");
    }
}

int process_single_submit(problem_runtime *pr_list) {
    int mem_lmt = 32;   // 32MB
    int time_lmt = 1000;    // 1000ms
    int topmemory = 0;

    int lang = LANGC;
    int i = 0;

    if(DEBUG) {
        printf("chdir to %s\n", pr_list->pathname);
    }
    chdir(pr_list->pathname);
    if(SUCCESS != compile(lang)) {
        update_solution(pr_list->oid, F_CE, 0, 0);
        return 0;
    }
    pid_t pidApp = fork();
    if(pidApp == 0) {
        char data_in[128];
        sprintf(data_in, "%s/%d/data.in", DATA_PATH, pr_list->pid);
        run_solution(1, 1, data_in);
    }
    else {
        int mem_peak, time_used;
        int jdg_ret;
        mem_peak = time_used = 0;
        watch_runtime(pidApp, LANGC, topmemory, mem_lmt, time_lmt, &mem_peak, &time_used);
        jdg_ret = judge_solution(pr_list->pid);
        if(DEBUG) {
            printf("Judge Result in Progress: %d\n", jdg_ret);
        }
        update_solution(pr_list->oid, jdg_ret, mem_peak, time_used);
        return 0;
    }
}

void process_all_submits(problem_runtime *pr_list[BATCH_SIZE]) {
    for(int i = 0; i < BATCH_SIZE; i++) {
        if(pr_list[i]->pathname[0] == '\0') continue;
        process_single_submit(pr_list[i]);
    }
}

void init_runtime(problem_runtime *pr_list[BATCH_SIZE]) {
    for(int i = 0; i < BATCH_SIZE; i++) {
        pr_list[i] = (problem_runtime*)malloc(sizeof(problem_runtime));
        pr_list[i]->pathname = (char*)malloc(sizeof(char) * 128);
        pr_list[i]->filename = (char*)malloc(sizeof(char) * 128);
        pr_list[i]->lang = (char*)malloc(sizeof(char) * 64);
        pr_list[i]->oid = (char*)malloc(sizeof(char) * 64);
        pr_list[i]->pid = -1;
    }
}

int main() {
    problem_runtime *pr_list[BATCH_SIZE];
    init_runtime(pr_list);
    process_smtqueue(pr_list);
    process_all_submits(pr_list);
}
