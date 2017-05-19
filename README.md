# Transparency Enforcing Cothority using Blockchains 

The Internet relies on many centralized services, e.g., for name resolution, authentication,
and content provisioning which provides authorities an easy way to monitor users or even censor unwanted content.
Even though the negative consequences are well-known, the centralization of the Internet has steadily increased even more over the last years.
One approach to mitigate these threats is to establish transparency mechanisms for authoritative records (such as DNS entries or TLS certificates) and expose them to public scrutiny.
Various systems, like Google's Certificate Transparency project or CONIKS provide transparency guarantees but only in a retroactive manner provided a victim has access to an honest monitor.
The goal of this project is to develop a proactive validation mechanism for record consistency to detect misbehavior before a client is deceived using blockchain technology and cothorities.
All implementations will be done with Google's Go programming language using the cothority framework of EPFL's DEDIS lab.