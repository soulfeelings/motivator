.PHONY: install dev build test clean

install:
	flutter pub get

dev:
	flutter run

build:
	flutter build apk --release
	flutter build ios --release --no-codesign

test:
	flutter test

clean:
	flutter clean
