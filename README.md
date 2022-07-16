# go-grpc

I've been seeing gRPC listed as a requirement in many job postings of recent, so I decided to do a little research on what gRPC was all about and how to use it in golang when I stumbled upon this amazingly detailed tutorial playlist from [TechSchool](https://www.youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf).

To understand what gRPC is, we must first have a basic understanding of RPC (Remote Procedure Call).

**Remote Procedure Call(RPC)** is when a computer program causes a procedure (subroutine) to execute in a different address space (commonly on another computer on a shared network), which is coded as if it were a normal (local) procedure call, without the programmer explicitly coding the details for the remote interaction. [Wikipedia](https://en.wikipedia.org/wiki/Remote_procedure_call)

> **Explanation to my 5 y/o self.** <br>
>
> Basically, you're able to call/execute a method of another program (usually running on another computer somwhere) , like it is a part of your program...... Cool, right? 😎

<br>

**`gRPC`**

---

In gRPC, a client application can directly call a method on a server application on a different machine as if it were a local object, making it easier for you to create distributed applications and services. As in many RPC systems, gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types. On the server side, the server implements this interface and runs a gRPC server to handle client calls. On the client side, the client has a stub (referred to as just a client in some languages) that provides the same methods as the server.

Read more on gRPC [here](https://grpc.io/docs/ "gRPC docs")

<br>
I'm pretty occupied at the moment, but I hope to make a side project of some sort that uses gRPC real soon.
