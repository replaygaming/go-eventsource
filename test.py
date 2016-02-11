#!/usr/bin/python
"""Creates, writes, and reads a labeled custom metric.

This is an example of how to use the Google Cloud Monitoring API to create,
write, and read a labeled custom metric.
The metric has two labels: color and size, and the data points represent
the number of shirts of the given color and size in inventory.

Prerequisites: Run this Python example on a Google Compute Engine virtual
machine instance that has been set up using these intructions:
https://cloud.google.com/monitoring/demos/setup_compute_instance.

Typical usage: Run the following shell commands on the instance:
    python write_labeled_metric.py --color yellow --size large --count 10
    python write_labeled_metric.py --color yellow --size medium --count 3
    python write_labeled_metric.py --color yellow --size large --count 12
    python write_labeled_metric.py --color blue --size medium --count 26
    python write_labeled_metric.py --color yellow --size large --count 8
    python write_labeled_metric.py --color blue --size medium --count 2
"""

import argparse
import time

from apiclient.discovery import build
import httplib2
from oauth2client.gce import AppAssertionCredentials

CUSTOM_METRIC_DOMAIN = "custom.cloudmonitoring.googleapis.com"
CUSTOM_METRIC_NAME = "%s/shirt_inventory" % CUSTOM_METRIC_DOMAIN


def GetProjectId():
  """Read the numeric project ID from metadata service."""
  http = httplib2.Http()
  resp, content = http.request(
      ("http://metadata.google.internal/"
       "computeMetadata/v1/project/numeric-project-id"),
      "GET", headers={"Metadata-Flavor": "Google"})
  if resp["status"] != "200":
    raise Exception("Unable to get project ID from metadata service.")
  return content


def GetNowRfc3339():
  """Retrieve the current time formatted per RFC 3339."""
  return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())


def CreateCustomMetric(service, project_id):
  """Create metric descriptor for the custom metric and send it to the API."""
  # You need to execute this operation only once. The operation is idempotent,
  # so, for simplicity, this sample code calls it each time

  # Create a label descriptor for each of the metric labels. The
  # "description" field should be more meaningful for your metrics.
  label_descriptors = []
  for label in ["color", "size",]:
    label_descriptors.append({"key": "/%s" % label,
                              "description": "The %s." % label})
  # Create the metric descriptor for the custom metric.
  metric_descriptor = {
      "name": CUSTOM_METRIC_NAME,
      "project": project_id,
      "typeDescriptor": {
          "metricType": "gauge",
          "valueType": "int64",
          },
      "labels": label_descriptors,
      "description": "The size of my shirt inventory.",
      }
  # Submit the custom metric creation request.
  try:
    request = service.metricDescriptors().create(
        project=project_id, body=metric_descriptor)
    request.execute()  # ignore the response
  except Exception as e:
    print "Failed to create custom metric: exception=%s" % e
    raise  # propagate exception


def WriteCustomMetric(service, project_id, now_rfc3339, color, size, count):
  """Write a data point to a single time series of the custom metric."""
  # Identify the particular time series to which to write the data by
  # specifying the metric and values for each label.
  timeseries_descriptor = {
      "project": project_id,
      "metric": CUSTOM_METRIC_NAME,
      "labels": {
          "%s/color" % CUSTOM_METRIC_DOMAIN: color,
          "%s/size" % CUSTOM_METRIC_DOMAIN: size,
          }
      }
  # Specify a new data point for the time series.
  timeseries_data = {
      "timeseriesDesc": timeseries_descriptor,
      "point": {
          "start": now_rfc3339,
          "end": now_rfc3339,
          "int64Value": count,
          }
      }
  # Submit the write request.
  request = service.timeseries().write(
      project=project_id, body={"timeseries": [timeseries_data,]})
  try:
    request.execute()   # ignore the response
  except Exception as e:
    print "Failed to write data to custom metric: exception=%s" % e
    raise  # propagate exception


def ReadCustomMetric(service, project_id, now_rfc3339, color, size):
  """Read all the timeseries data points for a given set of label values."""
  # To identify a time series, specify values for in label as a list.
  labels = ["%s/color==%s" % (CUSTOM_METRIC_DOMAIN, color),
            "%s/size==%s" % (CUSTOM_METRIC_DOMAIN, size),]
  # Submit the read request.
  request = service.timeseries().list(
      project=project_id,
      metric=CUSTOM_METRIC_NAME,
      youngest=now_rfc3339,
      labels=labels)
  # When a custom metric is created, it may take a few seconds
  # to propagate throughout the system. Retry a few times.
  start = time.time()
  while True:
    try:
      response = request.execute()
      for point in response["timeseries"][0]["points"]:
        print "  %s: %s" %  (point["end"], point["int64Value"])
      break
    except Exception as e:
      if time.time() < start + 20:
        print "Failed to read custom metric data, retrying..."
        time.sleep(3)
      else:
        print "Failed to read custom metric data, aborting: exception=%s" % e
        raise  # propagate exception


def main():
  # Define three parameters--color, size, count--each time you run the script"
  parser = argparse.ArgumentParser(description="Write a labeled custom metric.")
  parser.add_argument("--color", required=True)
  parser.add_argument("--size", required=True)
  parser.add_argument("--count", required=True)
  args = parser.parse_args()

  # Assign some values that will be used repeatedly.
  project_id = GetProjectId()
  now_rfc3339 = GetNowRfc3339()

  # Create a cloudmonitoring service object. Use OAuth2 credentials.
  credentials = AppAssertionCredentials(
      scope="https://www.googleapis.com/auth/monitoring")
  http = credentials.authorize(httplib2.Http())
  service = build(serviceName="cloudmonitoring", version="v2beta2", http=http)

  try:
    print "Labels: color: %s, size: %s." % (args.color, args.size)
    print "Creating custom metric..."
    CreateCustomMetric(service, project_id)
    time.sleep(2)
    print "Writing new data to custom metric timeseries..."
    WriteCustomMetric(service, project_id, now_rfc3339,
                      args.color, args.size, args.count)
    print "Reading data from custom metric timeseries..."
    ReadCustomMetric(service, project_id, now_rfc3339, args.color, args.size)
  except Exception as e:
    print "Failed to complete operations on custom metric: exception=%s" % e


if __name__ == "__main__":
  main()
