package lister

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
)

/**
  Utility component that iterates over all nodes in the memchached
  cluster and lists all their respective keys.
*/
type ClusterNodesKeyLister struct {
	clusterNodeUrls []string
}

//NewClusterNodeKeyLister - Returns a new ClusterNodeKeyLister instance.
func NewClusterNodeKeyLister(nodeUrls []string) *ClusterNodesKeyLister {

	return &ClusterNodesKeyLister{

		clusterNodeUrls: nodeUrls,
	}

}

//ListAllHostKeys - Lists all keys associated with all nodes in the cluster.
func (cnkl *ClusterNodesKeyLister) ListAllHostKeys() ([]string, error) {

	allClusterNodeKeys := make([]string, 1)

	for _, currentNode := range cnkl.clusterNodeUrls {

		currentNodeKeys, err := cnkl.listHostKeys(currentNode)

		if err != nil {

			log.Println("An error occured whilst attempting to fetch keys for node: " + currentNode + " " + err.Error())

		} else {

			allClusterNodeKeys = append(allClusterNodeKeys, currentNodeKeys...)
		}

	}

	return allClusterNodeKeys, nil
}

//private functions.

func (cnkl *ClusterNodesKeyLister) getNewConnection(hostAndPortNumber string) (net.Conn, error) {

	conn, err := net.Dial("tcp", hostAndPortNumber)

	if err != nil {

		return nil, err
	}

	return conn, nil

}

func (cnkl *ClusterNodesKeyLister) dispatchRequestAndReadResponse(connection net.Conn, command string, responseDelimiters []string) []string {
	fmt.Fprintf(connection, command)
	scanner := bufio.NewScanner(connection)
	var result []string

OUTER:
	for scanner.Scan() {
		line := scanner.Text()
		for _, delimeter := range responseDelimiters {
			if line == delimeter {
				break OUTER
			}
		}
		result = append(result, line)
		// if there is no delimiter specified, then the response is just a single line and we should return after
		// reading that first line (e.g. version command)
		if len(responseDelimiters) == 0 {
			break OUTER
		}
	}
	return result
}

func (cnkl *ClusterNodesKeyLister) listHostKeys(aHostAddressAndPort string) ([]string, error) {
	keys := []string{}
	conn, err := cnkl.getNewConnection(aHostAddressAndPort)
	if err != nil {

		log.Println("An error occured whilst attempting to connect to Memcached cluster node at: " + aHostAddressAndPort + " " + err.Error())

		return nil, err
	}

	//result := client.executer.execute("stats items\r\n", []string{"END"})
	result := cnkl.dispatchRequestAndReadResponse(conn, "stats items\r\n", []string{"END"})

	// identify all slabs and their number of items by parsing the 'stats items' command
	r, _ := regexp.Compile("STAT items:([0-9]*):number ([0-9]*)")
	slabCounts := map[int]int{}
	for _, stat := range result {
		matches := r.FindStringSubmatch(stat)
		if len(matches) == 3 {
			slabId, _ := strconv.Atoi(matches[1])
			slabItemCount, _ := strconv.Atoi(matches[2])
			slabCounts[slabId] = slabItemCount
		}
	}

	// For each slab, dump all items and add each key to the `keys` slice
	r, _ = regexp.Compile("ITEM (.*?) .*")
	for slabId, slabCount := range slabCounts {
		command := fmt.Sprintf("stats cachedump %v %v\n", slabId, slabCount)
		//commandResult := client.executer.execute(command, []string{"END"})
		commandResult := cnkl.dispatchRequestAndReadResponse(conn, command, []string{"END"})
		for _, item := range commandResult {
			matches := r.FindStringSubmatch(item)
			keys = append(keys, matches[1])
		}
	}

	conn.Close()

	return keys, nil
}
