import configparser
import os
import json
import datetime


parser = configparser.ConfigParser()
parser.read('test_config.ini')
configs = parser['test_config']

local_orb_path = configs.get("orb_path", os.path.dirname(os.getcwd()))
orb_address = configs.get("orb_address")

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


failing_scenarios = open(local_output_dir+"fail_scenarios.txt", "w")
failing_scenarios.write(str(datetime.datetime.now()))
failing_scenarios.close()

failing_scenarios = open(local_output_dir+"fail_scenarios.txt", "a")
failing_scenarios.write('\n\n')
failing_scenarios.write(f"ORB ADDRESS: {orb_address}")
failing_scenarios.write('\n\n')

for feature_index,  feature in enumerate(report['features']):
    if feature['status'] == "failed":
        for scenario_index, scenario in enumerate(feature['scenarios']):
            if scenario['status'] == "failed":
                for key, value in list(scenario.items()):
                    if key in ['name', 'duration', 'status', 'tags', 'filename', 'feature', 'error_msg']:
                        failing_scenarios.write('\n')
                        failing_scenarios.write(f"{key.upper()}: {value}")
                        print(f"{key.upper()}: {value}")
                    elif key == 'error_lines':
                        failing_scenarios.write('\n')
                        failing_scenarios.write(f"{key.upper()}:")
                        print(f"{key.upper()}:")
                        for line in scenario[key]:
                            line = line.split('\n')[0]
                            failing_scenarios.write('\n')
                            failing_scenarios.write(line)
                            print(line)
                    elif key == 'error_step':
                        failing_scenarios.write('\n')
                        failing_scenarios.write(f"{'step_type'.upper()} : {scenario[key]['step_type']}")
                        failing_scenarios.write('\n')
                        failing_scenarios.write(f"{'step_name'.upper()} : {scenario[key]['name']}")
                        print(f"{'step_type'.upper()} : {scenario[key]['step_type']}")
                        print(f"{'step_name'.upper()} : {scenario[key]['name']}")
                failing_scenarios.write('\n\n')
                failing_scenarios.write("---------------------------------------------------------------------------\n")
                print("-----------------------------------------------------------------------------------\n")

failing_scenarios.close()
