#!/bin/bash

# Create 8 players using the dart CLI
cd dart_client

dart run bin/mtg_cli.dart --token alice player signup alice
dart run bin/mtg_cli.dart --token bob player signup bob  
dart run bin/mtg_cli.dart --token charlie player signup charlie
dart run bin/mtg_cli.dart --token diana player signup diana
dart run bin/mtg_cli.dart --token eve player signup eve
dart run bin/mtg_cli.dart --token frank player signup frank
dart run bin/mtg_cli.dart --token grace player signup grace
dart run bin/mtg_cli.dart --token henry player signup henry
