from matplotlib.pyplot import polar
from numpy.core.fromnumeric import cumprod
from sklearn import svm

import os
import numpy as np
from matplotlib import pyplot as plt

class UVCharacteristics:
    date: str = None
    weatherMod: str = None
    MaxUV: int = None
    CurrentUV: int = None
    polarized: bool = None
    mean: float = None


def retreiveData(file, polarized=False):
    param = UVCharacteristics()
    with open(file, 'r') as f:
        content = f.read().splitlines()
        param.date = content[0]
        param.weatherMod = content[1].split(": ")[1]
        param.MaxUV = int(content[2].split(": ")[1])
        param.CurrentUV = int(content[3].split(": ")[1])
        param.polarized = polarized
        data = [int(i.split(",")[1]) for i in content[4::]]
    param.mean = (np.mean(data))
    return (param, data)

def getIndex(p: UVCharacteristics):
    # HARD CODED Voltage to INDEX Values UPDATE when more data is available
    def errorCalc(t):
        x,y = t
        return (abs(x-y),x)

    INDEXRANGES = {100:1, 1000:4, 400:2, 780:3, 40:0}
    if p.polarized:
        p.mean = p.mean*2.2 
    
    out = min((map(errorCalc, [(x,p.mean) for x in INDEXRANGES.keys()])))
    return INDEXRANGES[out[1]]

def buildTrainingSet(location):
    bigData = []
    labels = []
    directory = os.fsencode(location)
    for file in os.listdir(directory):
        if os.fsdecode(file).startswith('polarized'):
            indexParam, data = retreiveData(f"{location}/{os.fsdecode(file)}", polarized=True)
        else:
            indexParam, data = retreiveData(f"{location}/{os.fsdecode(file)}", polarized=False)
        
        bigData.append(data)
        labels.append(getIndex(indexParam))
    return (bigData, labels)

def buildPredictSet(fname: str, maxSize:int):
    with open(fname, 'r') as f:
        content = f.read().splitlines()
        data = [int(i.split(",")[1]) for i in content[4::]]
    if len(data) != maxSize:
        data = data[:maxSize]
    return data

def forceUniformity(data):
    maxDataSize = min(map(len,data))
    for i in range(len(data)):
        if len(data[i]) > maxDataSize:
            data[i] = data[i][:maxDataSize]
    return (data, maxDataSize)

def printData(UVIntensity, labels):
    for i in range(len(UVIntensity)):
        print(f"Average UV Index: {np.mean(UVIntensity[i])}")
        print(f"Current Label: {labels[i]}")
    plt.plot([np.mean(i) for i in UVIntensity], labels, 'bo')
    plt.title("UV Index vs Reading")
    plt.xlabel("GPIO IN")
    plt.ylabel("UV Index")
    plt.savefig('testing.png')

def execute_model():
    regr = svm.SVR()
    train_data, train_labels = buildTrainingSet("data")
    train_data_uniform, training_size = forceUniformity(train_data)
    printData(train_data, train_labels)

    regr.fit(train_data_uniform, train_labels)
    return (regr, training_size)

def predict(predictset, model) -> float:
    predict = model.predict([predictset])
    return predict

regr, maxSize = execute_model()

print(predict(buildPredictSet('data/predict_Sat_Mar_13_20-06-10_2021.txt', maxSize), regr))
