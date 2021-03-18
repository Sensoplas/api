from matplotlib import pyplot as plt
from sklearn import svm
from joblib import load

import os
import numpy as np

def buildPredictSet(inString: str, maxSize:int):
    newString = inString.split(',')
    data = [float(x) for x in newString]
    if len(data) > maxSize:
        data = data[::maxSize]
    else:
        meanData = np.mean(data)
        data = data + [meanData]*(maxSize-len(data))
    return data

def predict(predictset, model) -> float:
    predict = model.predict([predictset])
    return predict

import sys

if __name__ == "__main__":
    regr = load('/app/UVIndexModel.joblib') 
    #regr = load('UVIndexModel.joblib') 
    maxModelSize = int(os.environ["MODELSIZE"])
    #maxModelSize = 12076
    inString = sys.argv[1]
    predictSet = buildPredictSet(inString, maxModelSize)
    print(predict(predictSet, regr))
