import React, { PropTypes } from 'react';

class Image extends React.Component {
  constructor(props) {
    super(props);
    const src = this.props.src || this.props.errSrc;
    this.state = {
      src,
      isError: false,
    };
    this.handleError = this.handleError.bind(this);
  }

  handleError(event) {
    event.preventDefault();
    if (this.state.isError) {
      return;
    }
    this.setState({
      isError: true,
      src: this.props.errSrc,
    });
  }

  render() {
    return (
      <img {...this.props.imgProps}
        src={this.state.src}
        onError={this.handleError}
      />
    );
  }
}

Image.propTypes = {
  src: PropTypes.string.isRequired,
  errSrc: PropTypes.string.isRequired,
  imgProps: PropTypes.object,
};

export default Image;
