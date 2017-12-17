package main

import (
	"strconv"
)

//BlockSolver is a struct meant to help solve Blockchain hashes in a multi-threaded fashion
type BlockSolver struct {
	nworkers   int
	lastProof  int
	proofQueue chan int
	proofChan  chan int
	solved     bool
}

func (bs *BlockSolver) init(nworkers int, lastProof int) {
	bs.nworkers = nworkers
	bs.lastProof = lastProof
	bs.proofQueue = make(chan int, nworkers*2)
	bs.proofChan = make(chan int)
	bs.solved = false
}
func validProof(lastProof int, proof int) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(proof)
	guessHash := shaHash(guess)
	return guessHash[:3] == "000"
}

//TODO gotta stop workers once sync mutex indicates solution is found
func (bs *BlockSolver) proofOfWorker(workerID int) {
	for !bs.solved {
		proof := <-bs.proofQueue
		//fmt.Printf("Validating %d on worker %d \n", proof, workerID)
		if validProof(bs.lastProof, proof) {
			bs.solved = true
			bs.proofChan <- proof
		} else {
			bs.proofQueue <- proof + bs.nworkers
		}
	}
}
