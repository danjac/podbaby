import React from 'react';

class Image extends React.Component {
  constructor(props) {
    super(props);
    const src = this.props.src || this.props.errSrc;
    this.state = {
      src: src,
      isError: false
    };
  }

  handleError(event) {
    event.preventDefault()
    if (this.state.isError) {
      return;
    }
    this.setState({
      isError: true,
      src: this.props.errSrc
    });
  }

  render() {
    return (
      <img {...this.props.imgProps}
            src={this.state.src}
            onError={this.handleError.bind(this)} />
    );
  }
}

export default Image;
