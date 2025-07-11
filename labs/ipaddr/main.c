#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <arpa/inet.h>
#include <netinet/in.h>

struct foobar {
  int foo;
  char bar[16];
};

int foobarfunc() {

}

int main(int argc, char *argv[]){
  /*
    The addrinfo struct is used by network address and service translation functions like getaddrinfo().
    It contains the following fields:

    struct addrinfo {
      int              ai_flags;      // Input flags
      int              ai_family;     // Address family (AF_INET, AF_INET6, etc.)
      int              ai_socktype;   // Socket type (SOCK_STREAM, SOCK_DGRAM, etc.)
      int              ai_protocol;   // Protocol (IPPROTO_TCP, IPPROTO_UDP, etc.)
      socklen_t        ai_addrlen;    // Length of ai_addr
      struct sockaddr *ai_addr;       // Pointer to socket address structure
      char            *ai_canonname;  // Canonical name for hostname
      struct addrinfo *ai_next;       // Pointer to next addrinfo in linked list
    };
  */
  struct addrinfo hints, *res, *p;
  int status;
  char ipstr[INET6_ADDRSTRLEN];


  if (argc != 2) {
    fprintf(stderr, "usage: showip hostname\n");
    return 1;
  }

  // memset is a C library function that sets a block of memory to a specific value.
  // It is commonly used to initialize structures or arrays to zero or another value.
  // The function prototype is: void *memset(void *s, int c, size_t n);
  // It sets the first n bytes of the memory area pointed to by s to the byte value c.
  // For example, memset(&hints, 0, sizeof hints); sets all bytes of the 'hints' struct to 0.

  memset( &hints, 0, sizeof hints );

  hints.ai_family = AF_UNSPEC;
  hints.ai_socktype = SOCK_STREAM;


  if (( status = getaddrinfo(argv[1], NULL, &hints, &res )) != 0) {
    fprintf(stderr, "getaddrinfo %s\n", gai_strerror(status));
    return 2;
  }

  printf("IP addresses for %s:\n\n", argv[1]);

  for(p = res; p != NULL; p=p->ai_next) {
    void *addr;
    char *ipver;
    struct sockaddr_in *ipv4;
    struct sockaddr_in6 *ipv6;
    
    if (p->ai_family == AF_INET) {
      ipv4 = (struct sockaddr_in *)p->ai_addr;
      addr = &(ipv4->sin_addr);
      ipver = "IPv4";
    } else {
      ipv6 = (struct sockaddr_in6 *)p->ai_addr;
      addr = &(ipv6->sin6_addr);
      ipver = "IPV6";
    }

    inet_ntop(p->ai_family, addr, ipstr, sizeof ipstr);
    printf(" %s: %s\n", ipver, ipstr);
  }


  int sockfd;

  sockfd = socket(res->ai_family, res->ai_socktype, res->ai_protocol);

  int bind_ret = bind(sockfd, res->ai_addr, res->ai_addrlen);

  if (bind_ret == -1) {
    printf("Error bind number: %i", bind_ret);
  }

  freeaddrinfo(res);

  return 0;
}
