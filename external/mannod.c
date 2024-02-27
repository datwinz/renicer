#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]){
    if (argc == 1){
        printf("You should call a manpage\n");
        return(1);
    }
    char *man;
    char *arg;
    char *cmd;

    man = (char *) malloc(5);
    strcpy(man, "man ");
    // The longest manpagename I have is 74 chars, which is ridiculously long, so this is
    // probably fine.
    arg = (char *) malloc(75);
    strcpy(arg, argv[1]);

    cmd = (char *) malloc(80);
    strcpy(cmd, strcat(man, arg));
    system(cmd);

    free(man);
    free(arg);
    free(cmd);

    return(0);
}
