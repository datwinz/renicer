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

    man = (char *) malloc(10);
    strcpy(man, "sh -c \"man ");
    // The longest manpagename I have is 74 chars, which is ridiculously long, so this is
    // probably fine.
    arg = (char *) malloc(75);
    strcpy(arg, argv[1]);
    arg = strcat(arg, "\"");

    cmd = (char *) malloc(85);
    strcpy(cmd, strcat(man, arg));
    system(cmd);

    free(man);
    free(arg);
    free(cmd);

    return(0);
}
