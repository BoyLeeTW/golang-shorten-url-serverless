#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { GolangShortenUrlServerlessStack } from '../lib/golang-shorten-url-serverless-stack';

const app = new cdk.App();
new GolangShortenUrlServerlessStack(app, 'GolangShortenUrlServerlessStack');
