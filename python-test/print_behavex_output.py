import configparser
import os
import json


parser = configparser.ConfigParser()
parser.read('test_config.ini')
configs = parser['test_config']

local_orb_path = configs.get("orb_path", os.path.dirname(os.getcwd()))

local_output_dir = local_orb_path + "/python-test/output/"

with open(local_output_dir + "overall_status.json") as f:
    print("-----------------------------------------------------------------------------------\n")
    suite_result = json.load(f)['status']
    print(f"Test suite result: {suite_result}\n")
    print("-----------------------------------------------------------------------------------\n")
    f.close()


file = open(local_output_dir+"report.json")

report = json.load(file)

file.close()

failed_scenarios = dict()

for feature_index,  feature in enumerate(report['features']):
    if feature['status'] == "failed":
        for scenario_index, scenario in enumerate(feature['scenarios']):
            if scenario['status'] == "failed":
                for key, value in list(scenario.items()):
                    print(f"{key.upper()}: {value}")
                    failed_scenarios.update(report['features'][feature_index]['scenarios'][scenario_index])
                print("-----------------------------------------------------------------------------------\n")
with open(f"{local_output_dir}failed_scenarios.json", "w") as f:
    json.dump(failed_scenarios, f)

f.close()
